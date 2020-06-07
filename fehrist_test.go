package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kadnan/fehrist/fehrist"
)

func TestDocumentFileNotFoundCSV(t *testing.T) {
	path := "/Users/AdnanAhmad/Data/Development/PetProjects/Fehrist/"
	fileName := path + "LOL.csv"
	CSVDocument := &fehrist.CSV{IndexName: "local"}
	indexCount, _ := CSVDocument.Index(fileName)

	if indexCount != -1 {
		t.Errorf("Test DocumentFileNotFound Failed")
	}
}

func TestDocumentIndexedCSV(t *testing.T) {
	path := "/Users/AdnanAhmad/Data/Development/PetProjects/Fehrist/"
	fileName := path + "1.csv"
	CSVDocument := &fehrist.CSV{IndexName: "local"}
	indexCount, _ := CSVDocument.Index(fileName)
	if indexCount < 1 {
		t.Errorf("Test Document Index Failed CSV")
	}
}

func TestDocumentFileNotFoundJSON(t *testing.T) {
	path := "/Users/AdnanAhmad/Data/Development/PetProjects/Fehrist/"
	fileName := path + "LOL.json"
	CSVDocument := &fehrist.CSV{IndexName: "local"}
	indexCount, _ := CSVDocument.Index(fileName)
	if indexCount != -1 {
		t.Errorf("Test Document Indexed Failed")
	}
}

func TestDocumentIndexedJSON(t *testing.T) {
	path := "/Users/AdnanAhmad/Data/Development/PetProjects/Fehrist/"
	fileName := path + "1.json"
	JSONDocument := &fehrist.JSON{IndexName: "local"}
	indexCount, _ := JSONDocument.Index(fileName)
	if indexCount < 1 {
		t.Errorf("Test Document Index Failed JSON")
	}
}
func TestInitCSV(t *testing.T) {
	Document := &fehrist.CSV{IndexName: "local"}
	result := Document.Init()
	if result != 1 {
		t.Errorf("Test Initialization failed for CSV Document")
	}
}
func TestInitJSON(t *testing.T) {
	Document := &fehrist.JSON{IndexName: "local"}
	result := Document.Init()
	if result != 1 {
		t.Errorf("Test Initialization failed for JSON Document")
	}
}

func TestSearchCSVKWExist(t *testing.T) {
	var input map[string]interface{}

	Document := &fehrist.CSV{IndexName: "local"}
	result := Document.Init()

	if result == 1 {
		out, _, _ := Document.Search("mango")
		err := json.Unmarshal([]byte(out), &input)
		if err != nil {
			fmt.Println(err)
		}
		_, ok := input["Total"]
		if !ok {
			t.Errorf("Test TestSearchCSVKWExist failed for CSV Document")
		}
	}
}

func TestSearchJSONKWExist(t *testing.T) {
	var input map[string]interface{}

	Document := &fehrist.JSON{IndexName: "local"}
	result := Document.Init()

	if result == 1 {
		out, _, _ := Document.Search("mango")
		err := json.Unmarshal([]byte(out), &input)
		if err != nil {
			fmt.Println(err)
		}
		_, ok := input["Total"]
		if !ok {
			t.Errorf("Test TestSearchCSVKWExist failed for JSON  Document")
		}
	}
}
