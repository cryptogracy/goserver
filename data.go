package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type dataControl interface {
	AddFile(hash string, lifespan int) error
  AddMeta(Id string, lifespan int) error
	Cleanup() (int64, error)
	Check(hash string) bool
	Close() error
}

type DB struct {
	db *sql.DB
}

func DBInit(database string) (DB, error) {
	db, err := sql.Open("sqlite3", database)
	ret := DB{db}
	if err != nil {
		return ret, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Files
		(
			Hash     STRING    PRIMARY KEY,
			Death    TIMESTAMP NOT NULL,
			Short    STRING
		)`)
	if err != nil {
		return DB{db}, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Meta
		(
			Id    STRING PRIMARY KEY,
			DEATH TIMESTAMP NOT NULL
		)`)
	return DB{db}, err
}

func (db DB) Close() error { return db.db.Close() }

func (db DB) AddFile(hash string, lifespan int) error {
	entry := struct {
		Hash  string
		Death time.Time
		Short string
	}{hash,
		time.Now().Add(time.Duration(lifespan) * time.Second),
		""}
	_, err := db.db.Exec(
    "INSERT INTO Files (Hash, Death, Short) VALUES (?, ?, ?)",
		entry.Hash, entry.Death, entry.Short)
	return err
}

func (db DB) AddMeta(id string, lifespan int) error {
	entry := struct {
		Id    string
		Death time.Time
	}{id, time.Now().Add(time.Duration(lifespan) * time.Second)}
	_, err := db.db.Exec("INSERT INTO Files (Hash, Death) VALUES (?, ?)",
		entry.Id, entry.Death)
	return err

}

func (db DB) Cleanup() (int64, error) {
	result, err := db.db.Exec(
		"DELETE FROM Files WHERE death < $1", time.Now())
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (db DB) Check(hash string) bool {
	var res string
	if err := db.db.QueryRow("SELECT Hash from Files WHERE Hash = $1", hash).Scan(&res); err == sql.ErrNoRows {
		return false
	} else if err != nil {
		panic(err)
	}
	return true
}
