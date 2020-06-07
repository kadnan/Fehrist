// Package fehrist implements routines related to different kind of indexing

package fehrist

import (
	"bytes"
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack"
)

var entries = make(map[string][2]string) //to store values
var tokenized = make(map[string]string)
var mergedIndexMap = make(map[string]string)
var mergedMap = make(map[string]string)
var mergedMapDocuments = make(map[string][2]string)
var idx = 0
var isLoaded = false

//All Constants
const success int = 1
const failure int = -1

/* unmarshalData removes all the field that are not required and returns a flat Map
Author: Peter Hellberg (https://gophers.slack.com/) - Thanks Peter!!
*/
func unmarshalData(data []byte) ([]map[string]interface{}, error) {
	var input []map[string]interface{}

	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}

	var resp []map[string]interface{}

	for _, in := range input {
		d := map[string]interface{}{}
		for k, v := range in {
			switch v.(type) {
			case string, float64:
				d[k] = v
			}
		}
		resp = append(resp, d)
	}

	return resp, nil
}

// Save function saves the value in the in a file.
func save(content map[string][2]string, t map[string]string, fileHandleDocument *os.File, fileHandleIndex *os.File) (int8, error) {
	// if !isLoaded {
	// 	return -1, errors.New("File not found")
	// }

	defer fileHandleDocument.Close()
	defer fileHandleIndex.Close()

	// Saving Content Data after Marshalizing
	if len(content) > 0 {
		b, err := msgpack.Marshal(content)

		if err != nil {
			return 0, errors.New("Decoding Failed for Document")
		}

		fileHandleDocument.Write(b)
	}

	// Saving Index Data after Marshalizing
	if len(t) > 0 {
		b, err := msgpack.Marshal(t)

		if err != nil {
			return 0, errors.New("Decoding Failed for Index")
		}

		fileHandleIndex.Write(b)

	}

	return 1, nil
}

/* saveIndex saves all index and document related info on disk. It is responsible for:
- Create a Folder of Index Name
- For each document index it creates a file with numeric sequence name with extension .idx
- It stores original document along with assignedID
*/
func saveIndex(indexName string, path string, documentFileName string) {
	fileSquence := 0
	indexPath := indexName + "/"
	pattern := filepath.Join(indexPath, "*.idx")

	if _, err := os.Stat(indexName); os.IsNotExist(err) {
		os.Mkdir(indexName, 0700) //Write from the same program
	}

	// Folder Created. Now we have to check the next available sequence of file

	existingIndexFiles, err := filepath.Glob(pattern)

	fileSquence = len(existingIndexFiles)

	if err != nil {
		fmt.Println(err.Error())
	}

	if len(existingIndexFiles) > 0 {
		fileSquence = len(existingIndexFiles)
	}

	indexFileName := indexPath + strconv.Itoa(fileSquence) + ".idx"
	docFileName := indexPath + documentFileName + ".document"

	// Save the Document File with .document extension
	docFile, err := os.OpenFile(docFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}

	indexFile, err := os.OpenFile(indexFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = save(entries, tokenized, docFile, indexFile)

	if err != nil {
		fmt.Println(err.Error())
	}
}

// generateDocID generates a random DocID
func generateDocID(text string) string {
	algorithm := sha1.New()
	algorithm.Write([]byte(text))
	idx++
	return hex.EncodeToString(algorithm.Sum(nil)) //ALERT: Implement hex based ID
	//return strconv.Itoa(idx)
}

// Indexer is the interface that implements important stuff
type Indexer interface {
	Index(fileName string)
	assignDocID(entry string)
	tokenizeDocument() map[string]string
}

// A CSV represents a CSV Doccument
type CSV struct {
	IndexName string
}

// A JSON represents a CSV Doccument
type JSON struct {
	IndexName string
}

//DocumentList holds the return searched doc structure
type DocumentList struct {
	FileName string
	DocText  string
}

//SearchResult implements Search JSON
type SearchResult struct {
	Total  int
	Result []DocumentList
}

func msgPack2MapIndex(marshalled string) map[string]string {
	var tempTokenizedMap = make(map[string]string)
	msgpack.Unmarshal([]byte(marshalled), &tempTokenizedMap)
	return tempTokenizedMap
}
func msgPack2MapDocument(marshalled string) map[string][2]string {
	var tempDocMap = make(map[string][2]string)
	msgpack.Unmarshal([]byte(marshalled), &tempDocMap)
	return tempDocMap
}

//generateJSONArray checks whether the JSON is array of object or not, if no then make it one
func generateJSONArray(data string) []byte {
	//jsonWithoutSpace := strings.ReplaceAll(string(data), " ", "")
	jsonWithoutSpace := strings.TrimSpace(data)
	if string(jsonWithoutSpace[0]) != "[" && string(jsonWithoutSpace[len(jsonWithoutSpace)-1]) != "]" {
		jsonWithoutSpace = "[" + jsonWithoutSpace + "]"
	}

	return []byte(jsonWithoutSpace)
}

//Index is used to index JSON documents after assigning Document ID
func (c *JSON) Index(fileName string) (int, error) {
	_, fileNameOnly := filepath.Split(fileName)
	// Read the file
	file, _ := os.Open(fileName)
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	s := string(byteValue)
	fixedJSON := generateJSONArray(s)

	cleanJSON, err := unmarshalData(fixedJSON)
	for _, v := range cleanJSON {
		b, _ := json.Marshal(v)
		c.assignDocID(string(b), fileNameOnly)
	}

	c.tokenizeDocument()
	// Index is created now save the Index files and original mapped document in files
	saveIndex(c.IndexName, ".", fileNameOnly)

	if err != nil {
		return failure, nil
	}
	return len(tokenized), nil
}

func (c *JSON) assignDocID(entry string, documentFile string) {
	var rec [2]string
	rec[0] = documentFile
	rec[1] = entry
	now := time.Now()
	docID := "f_" + generateDocID(now.String())
	entries[docID] = rec
}

//tokenzeDocument tokenize the document into words and store them into Array.
func (c *JSON) tokenizeDocument() {

	rec := make(map[string]string)
	for key, entry := range entries {
		key = strings.TrimSpace(key)

		json.Unmarshal([]byte(entry[1]), &rec)
		for _, v := range rec {

			val := strings.ToLower(v)
			words := strings.Fields(val)
			for _, word := range words {
				_, ok := tokenized[word]

				if ok {
					tokenized[word] = tokenized[word] + "|" + key
				} else {
					tokenized[word] = key
				}
			}
		}
	}

}

//Init initializes the index and document related maps of the given index for JSON Documents
func (c *JSON) Init() int {
	var tempTokenizedMap = make(map[string]string)
	var tempDocmentMap = make(map[string][2]string)
	path, err := os.Getwd()

	//Fetching all index files and merge their maps into a single map
	indexPath := c.IndexName + "/"
	pattern := filepath.Join(indexPath, "*.idx")
	existingIndexFiles, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println(err.Error())
	}
	pattern = filepath.Join(indexPath, "*.document")
	existingDocumentFiles, err := filepath.Glob(pattern)

	if err != nil {
		fmt.Println(err.Error())
	}

	// Iterating Index files and map merging.
	for _, z := range existingIndexFiles {
		file, _ := os.Open(filepath.Join(path, z))
		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		contents := buf.String()
		file.Close() // Close the file immediately once done
		tempTokenizedMap = msgPack2MapIndex(contents)

		if len(tempTokenizedMap) > 0 {
			for key, value := range tempTokenizedMap {

				if _, found := mergedMap[key]; found {
					mergedMap[key] = mergedMap[key] + "|" + value
				} else {
					mergedMap[key] = value
				}

			}
		}
	}

	//Iterating Document files and merge them
	for _, z := range existingDocumentFiles {
		file, _ := os.Open(filepath.Join(path, z))

		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		contents := buf.String()
		file.Close() // Close the file immediately once done
		tempDocmentMap = msgPack2MapDocument(contents)

		if len(tempDocmentMap) > 0 {
			for key, value := range tempDocmentMap {
				mergedMapDocuments[key] = value
			}
		}
	}
	if err != nil {
		return -1
	} else {
		return 1
	}
}

// Search returns the result against the keyword being provided.
func (c *JSON) Search(keyword string) (string, int, error) {
	var documents []DocumentList
	// var result SearchResult
	//var docs []string

	if len(mergedMap) == 0 && len(mergedMapDocuments) == 0 {
		return "", -1, errors.New("No data was found. Did you call Init function?")
	}

	//Check the index map first
	v, found := mergedMap[strings.ToLower(keyword)]
	keys := strings.Split(v, "|")

	//result = SearchResult{Total:len(keys),Result: }

	//fmt.Println(mergedMapDocuments[keys[0]])

	for _, documentID := range keys {
		entry := mergedMapDocuments[documentID]

		if len(entry) == 2 {
			documents = append(documents, DocumentList{FileName: entry[0], DocText: entry[1]})
		}
	}
	x := SearchResult{Total: len(keys), Result: documents}
	jsonData, err := json.Marshal(x)
	if err != nil {
		return "", -1, errors.New("Could not decode")
	}

	if found {
		return string(jsonData), 1, nil
	}
	return "", -1, nil
}

//Init initializes the index and document related maps of the given index for CSV Documents
func (c *CSV) Init() int {
	var tempTokenizedMap = make(map[string]string)
	var tempDocmentMap = make(map[string][2]string)
	path, _ := os.Getwd()

	//Fetching all index files and merge their maps into a single map
	indexPath := c.IndexName + "/"
	pattern := filepath.Join(indexPath, "*.idx")
	existingIndexFiles, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println(err.Error())
	}
	pattern = filepath.Join(indexPath, "*.document")
	existingDocumentFiles, err := filepath.Glob(pattern)

	if err != nil {
		fmt.Println(err.Error())
	}

	// Iterating Index files and map merging.
	for _, z := range existingIndexFiles {
		file, _ := os.Open(filepath.Join(path, z))
		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		contents := buf.String()
		file.Close() // Close the file immediately once done
		tempTokenizedMap = msgPack2MapIndex(contents)

		if len(tempTokenizedMap) > 0 {
			for key, value := range tempTokenizedMap {

				if _, found := mergedMap[key]; found {
					mergedMap[key] = mergedMap[key] + "|" + value
				} else {
					mergedMap[key] = value
				}

			}
		}
	}

	//Iterating Document files and merge them
	for _, z := range existingDocumentFiles {
		file, _ := os.Open(filepath.Join(path, z))

		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		contents := buf.String()
		file.Close() // Close the file immediately once done
		tempDocmentMap = msgPack2MapDocument(contents)

		if len(tempDocmentMap) > 0 {
			for key, value := range tempDocmentMap {
				mergedMapDocuments[key] = value
			}
		}
	}
	if err != nil {
		return -1
	} else {
		return 1
	}
}

// Search returns the result against the keyword being provided.
func (c *CSV) Search(keyword string) (string, int, error) {
	var documents []DocumentList
	// var result SearchResult
	//var docs []string

	if len(mergedMap) == 0 && len(mergedMapDocuments) == 0 {
		return "", -1, errors.New("No data was found. Did you call Init function?")
	}

	//Check the index map first
	v, found := mergedMap[strings.ToLower(keyword)]
	keys := strings.Split(v, "|")

	for _, documentID := range keys {
		entry := mergedMapDocuments[documentID]

		if len(entry) == 2 {
			documents = append(documents, DocumentList{FileName: entry[0], DocText: entry[1]})
		}
	}
	x := SearchResult{Total: len(keys), Result: documents}
	jsonData, err := json.Marshal(x)
	if err != nil {
		return "", -1, errors.New("Could not decode")
	}

	if found {
		return string(jsonData), 1, nil
	}
	return "", 1, nil
}

func (c *CSV) assignDocID(entry string, documentFile string) {
	var rec [2]string
	rec[0] = documentFile
	rec[1] = entry
	now := time.Now()
	docID := "f_" + generateDocID(now.String())
	entries[docID] = rec
}

//tokenzeDocument tokenize the document into words and store them into Array.
func (c *CSV) tokenizeDocument() {

	for key, entry := range entries {
		key = strings.TrimSpace(key)
		line := strings.Replace(entry[1], ",", " ", 3)
		line = strings.ToLower(line)
		words := strings.Fields(line)

		for _, word := range words {
			_, ok := tokenized[word]
			if ok {
				tokenized[word] = tokenized[word] + "|" + key
			} else {
				tokenized[word] = key
			}
		}
	}
}

//Index indexes the document
func (c *CSV) Index(fileName string) (int, error) {

	_, fileNameOnly := filepath.Split(fileName)
	// Read the file
	file, _ := os.Open(fileName)
	defer file.Close()

	parser := csv.NewReader(file)
	parser.FieldsPerRecord = -1

	if _, err := parser.Read(); err != nil {
		return failure, errors.New("File not found")
	}

	records, err := parser.ReadAll()

	if err != nil {
		return failure, errors.New("Could not read the CSV file")
	}

	// Assign Document ID to each record
	for _, record := range records {
		rec := strings.Join(record, ",")
		c.assignDocID(rec, fileNameOnly)
	}
	if len(entries) > 0 {
		c.tokenizeDocument()
		// Index is created now save the Index files and original mapped document in files
		saveIndex(c.IndexName, ".", fileNameOnly)
	}

	return len(tokenized), nil
}
