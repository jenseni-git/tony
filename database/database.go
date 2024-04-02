package database

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(file string) *sql.DB {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.WithField("src", "database").WithError(err).Fatal("Failed to open database")
	}

	return db
}
