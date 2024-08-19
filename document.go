package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

type document struct {
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Text  string `xml:"abstract"`
	ID    int
}

func loadDocuments(filename string) ([]document, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dump := struct {
		Documents []document `xml:"doc"`
	}{}

	start := time.Now()
	dec := xml.NewDecoder(f)
	if err := dec.Decode(&dump); err != nil {
		return nil, err
	}
	fmt.Println("time to decode the xml file", time.Since(start))

	docs := dump.Documents
	for i := range docs {
		docs[i].ID = i
	}

	return docs, nil
}
