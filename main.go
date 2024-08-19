package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/adilsitos/fts/persistency"
)

func main() {
	// forceRead := flag.Bool("forceRead", false, "force to read the wiki data file")
	deleteBackup := flag.Bool("recreateBackup", false, "recrete the index backup")
	flag.Parse()

	docs, err := loadDocuments("./enwiki-latest-abstract1.xml")
	if err != nil {
		log.Panic(err)
	}

	idx := make(index)
	idx.add(docs)

	err = createBackup(idx, *deleteBackup)
	if err != nil {
		log.Panic(err)
	}

	res := idx.search("small wild cat")
	for _, i := range res {
		fmt.Println(docs[i])
	}

}

func createBackup(idx index, deleteBkpFile bool) error {
	var err error
	backupFile := "backup.txt"

	if deleteBkpFile {
		os.Remove(backupFile)
	}

	engine, err := persistency.NewEngine(backupFile)
	if err != nil {
		return err
	}

	for key, val := range idx {
		err = engine.Set(key, val)
		if err != nil {
			return err
		}
	}

	return nil
}
