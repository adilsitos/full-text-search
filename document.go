package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type document struct {
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Text  string `xml:"abstract"`
	ID    int
}

func loadDocuments(path string) ([]document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := xml.NewDecoder(f)
	dump := struct {
		Documents []document `xml:"doc"`
	}{}

	if err := dec.Decode(&dump); err != nil {
		return nil, err
	}

	docs := dump.Documents
	for i := range docs {
		docs[i].ID = i
	}

	fmt.Println("read document")

	return docs, nil
}
