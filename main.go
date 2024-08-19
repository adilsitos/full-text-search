package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/adilsitos/fts/persistency"
)

func main() {
	forceRead := flag.Bool("forceRead", false, "force to read the wiki data file")
	deleteBackup := flag.Bool("recreateBackup", false, "recrete the index backup")
	flag.Parse()

	idx, err := createIdx("./enwiki-latest-abstract1.xml", *forceRead)
	if err != nil {
		log.Panic(err)
	}

	err = createBackup(idx, *deleteBackup)
	if err != nil {
		log.Panic(err)
	}

	// res := idx.search("small wild cat")
	// for _, i := range res {
	// 	fmt.Println(docs[i])
	// }
}

func createIdx(filename string, forceRead bool) (index, error) {
	engine, err := persistency.NewEngine("backup.txt")
	if err != nil {
		return nil, err
	}

	_, fileStatErr := os.Stat("backup.txt")

	if !errors.Is(fileStatErr, os.ErrNotExist) && !forceRead {
		_, backup := engine.GetMapFromFile()
		return backup, nil
	}

	docs, err := loadDocuments(filename)
	if err != nil {
		return nil, err
	}

	idx := make(index)
	idx.add(docs)

	return idx, nil
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

	engine.Restore()

	return nil
}
