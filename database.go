package loctogene

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func GetDB(file string) (*sql.DB, error) {
	fmt.Printf("Opening db %s...\n", file)

	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, err
	}

	return db, nil
}
