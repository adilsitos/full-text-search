package main

import (
	"database/sql"
	"fmt"
	"os"
)

func ConnectIntoDB() *sql.DB {
	dbName := "file:./local.db"

	db, err := sql.Open("libsql", dbName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}

	return db
}
