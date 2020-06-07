package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kadnan/fehrist/fehrist"
)

func main() {
	path, _ := os.Getwd()

	//Indexing CSV Files
	CSVDocument := &fehrist.CSV{IndexName: "local"}
	for i := 1; i < 3; i++ {
		fileName := path + "/" + strconv.Itoa(i) + ".csv"
		fmt.Println("Indexing CSV data from the file,", fileName, ". Please wait...")

		indexCount, err := CSVDocument.Index(fileName)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Total Words indexed", indexCount)
		}
	}

	//Indexing JSON files
	JSONDocument := &fehrist.JSON{IndexName: "local"}
	for i := 1; i < 3; i++ {
		fileName := path + "/" + strconv.Itoa(i) + ".json"
		fmt.Println("Indexing CSV data from the file,", fileName, ". Please wait...")

		indexCount, err := JSONDocument.Index(fileName)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Total Words indexed", indexCount)
		}
	}

	/* Searching Documents */

	CSVDocument.Init()
	result, _, err := CSVDocument.Search("siddiqi")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Printing the text present in CSV Document")
	fmt.Println(result)

	JSONDocument.Init()
	result, searchCount, err := JSONDocument.Search("mango")
	fmt.Println(searchCount)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Printing the text present in JSON Document")
	fmt.Println(result)

}
