package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type dataControl interface {
	Add(hash string, lifespan int) error
	Cleanup() error
  Check(hash string) bool
}

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

func DBInit() (DB, error) {
	db, err := sql.Open("sqlite3", configuration.Database)
	ret := DB{db}
	if err != nil {
		return ret, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Files
		(
			Hash     STRING    PRIMARY KEY,
			Birth    TIMESTAMP NOT NULL,
			Death    TIMESTAMP NOT NULL,
			Lifetime TIMESTAMP NOT NULL,
			Short    STRING
		)`)
	return DB{db}, err
}

func (db DB) Add(hash string, lifespan int) error {
	now := time.Now()
	dur := time.Duration(lifespan) * time.Second
	entry := FileEntry{hash, now, now.Add(dur), dur, ""}
	_, err := db.db.Exec(
		"INSERT INTO Files (Hash, Birth, Death, Lifetime) VALUES (?, ?, ?, ?)",
		entry.Hash, entry.Birth, entry.Death, entry.Lifetime)
	return err
}

func (db DB) Cleanup() error {
	result, err := db.db.Exec(
		"DELETE FROM Files WHERE death < $1", time.Now())
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Deleted %v old Files", affected)
	return nil
}

func (db DB) Check(hash string) bool {
	var res string
	err := db.db.QueryRow("SELECT Hash from Files WHERE Hash = $1", hash).Scan(&res)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		errorPanic(err)
	}
	return true
}
