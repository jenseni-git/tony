package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(file string) *sql.DB {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}

	return db
}
