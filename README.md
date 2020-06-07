# Fehrist
[![Build Status](https://api.travis-ci.org/kadnan/fehrist.svg)](https://travis-ci.org/kadnan/fehrist)

_Fehrist_ is a pure Go library for indexing different types of documents. Currently it supports only CSV and JSON but flexible architecture gives you liberty to add more documents. Fehrist(فہرست) is an Urdu word for **Index**. Similar terminologies used in Arabic(فھرس) and Farsi(فہرست) as well.

Fehrist is based on [Inverted Index]([https://en.wikipedia.org/wiki/Inverted_index]) data structure for indexing purposes.

## Examples
### For indexing
```
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
}
```
If you want to learn how this all work then visit the [blog post](http://blog.adnansiddiqi.me/fehrist-document-indexing-library-in-go/)