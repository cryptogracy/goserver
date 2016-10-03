package main

import (
	"database/sql"
	"fmt"
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
	Short    string
}

func (db *DB) Init() {
	var err error
	db.db, err = sql.Open("sqlite3", configuration.Database)
	errorPanic(err)

	_, err = db.db.Exec(`
		CREATE TABLE IF NOT EXISTS Files
		(
			Hash     STRING    PRIMARY KEY,
			Birth    TIMESTAMP NOT NULL,
			Death    TIMESTAMP NOT NULL,
			Lifetime TIMESTAMP NOT NULL,
			Short    STRING
		)`)
	errorPanic(err)
}

func (db *DB) AddFile(hash string, lifespan int) {
	now := time.Now()
	dur := time.Duration(lifespan) * time.Second
	db.addFile(FileEntry{hash, now, now.Add(dur), dur, ""})
}

func (db *DB) addFile(entry FileEntry) {
	_, err := db.db.Exec(
		"INSERT INTO Files (Hash, Birth, Death, Lifetime) VALUES (?, ?, ?, ?)",
		entry.Hash, entry.Birth, entry.Death, entry.Lifetime)
	errorPanic(err)
}

func (db *DB) RemoveOldFile() {
	result, err := db.db.Exec(
		"DELETE FROM Files WHERE death < $1", time.Now())
	errorPanic(err)
	affected, err := result.RowsAffected()
	errorPanic(err)
	log.Printf("Deleted %v old Files", affected)
}

func (db *DB) RemoveFilePeriod(period time.Duration) {
	for true {
		db.RemoveOldFile()
		time.Sleep(period)
	}
}

func (db *DB) CheckPresence(hash string) bool {
	var res string
	err := db.db.QueryRow("SELECT Hash from Files WHERE Hash = $1", hash).Scan(&res)
	fmt.Println(res)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		errorPanic(err)
	}
	return true
}
