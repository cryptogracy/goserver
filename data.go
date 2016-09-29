package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type DB struct {
	db *sql.DB
}

type FileEntry struct {
	Hash     string
	Birth    time.Time
	Death    time.Time
	Lifetime time.Duration
}

func (db *DB) Init() {
	var err error
	db.db, err = sql.Open("sqlite3", configuration.Database)
	stdErrCheck(err)

	_, err = db.db.Exec(`
		CREATE TABLE IF NOT EXISTS Files
		(
			Hash     STRING    PRIMARY KEY,
			Birth    TIMESTAMP NOT NULL,
			Death    TIMESTAMP NOT NULL,
			Lifetime TIMESTAMP NOT NULL
		)`)
	stdErrCheck(err)
}

func (db *DB) AddFile(entry FileEntry) {
	_, err := db.db.Exec(
		"INSERT INTO Files (Hash, Birth, Death, Lifetime) VALUES (?, ?, ?, ?)",
		entry.Hash, entry.Birth, entry.Death, entry.Lifetime)
	if err != nil {
		log.Println(err)
	}
	stdErrCheck(err)
}
