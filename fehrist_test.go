package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/kadnan/fehrist/fehrist"
)

func TestDocumentFileNotFoundCSV(t *testing.T) {
	path, _ := os.Getwd()
	fileName := path + "/" + "LOL.csv"
	CSVDocument := &fehrist.CSV{IndexName: "local"}
	indexCount, _ := CSVDocument.Index(fileName)

	if indexCount != -1 {
		t.Errorf("Test DocumentFileNotFound Failed")
	}
}

func TestDocumentIndexedCSV(t *testing.T) {
	path, _ := os.Getwd()
	fileName := path + "/" + "1.csv"
	CSVDocument := &fehrist.CSV{IndexName: "local"}
	indexCount, err := CSVDocument.Index(fileName)
	fmt.Println(indexCount)
	fmt.Println(err)

	if indexCount < 1 {
		t.Errorf("Test Document Index Failed CSV")
	}
}

func TestDocumentFileNotFoundJSON(t *testing.T) {
	path, _ := os.Getwd()
	fileName := path + "/" + "LOL.json"
	CSVDocument := &fehrist.CSV{IndexName: "local"}
	indexCount, _ := CSVDocument.Index(fileName)
	if indexCount != -1 {
		t.Errorf("Test Document Indexed Failed")
	}
}

func TestDocumentIndexedJSON(t *testing.T) {
	path, _ := os.Getwd()
	fileName := path + "/" + "1.json"
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
