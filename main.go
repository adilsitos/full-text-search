package main

import (
	"fmt"
	"log"
)

func main() {
	docs, err := loadDocuments("./enwiki-latest-abstract1.xml")
	if err != nil {
		log.Panic(err)
	}

	idx := make(index)
	idx.add(docs)

	// idx.add([]document{{ID: 1, Text: "A donut on a glass plate. Only the donuts."}})
	// idx.add([]document{{ID: 2, Text: "donut is a donut"}})

	res := idx.search("Small wild cat")

	for _, i := range res {
		fmt.Println(docs[i])
	}
}
