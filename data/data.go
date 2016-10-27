package data

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DB struct {
	db *gorm.DB
}

type File struct {
	gorm.Model
	Hash  string    `gorm:"primary_key;unique"`
	Death time.Time `gorm:"not null"`
	Short string
}

type Metadata struct {
	gorm.Model
	Files []File
	Short string    `gorm:"not null;unique"`
	Death time.Time `gorm:"not null"`
}

func Init(database string) (DB, error) {

	rand.Seed(time.Now().UTC().UnixNano())
	db := DB{}
	var err error
	db.db, err = gorm.Open("sqlite3", database)
	if err != nil {
		return db, err
	}

	db.db.LogMode(false)

	db.db.AutoMigrate(&File{}, &Metadata{})

	return db, nil
}

func (db DB) Close() error {
	return db.db.Close()
}

func (db DB) AddFile(hash string, lifespan int) bool {
	death := time.Now().Add(time.Duration(lifespan) * time.Second)
	short := ""

	file := File{Hash: hash, Death: death, Short: short}
	db.db.Create(&file)

	if db.db.NewRecord(file) {
		return false
	}
	return true
}

func (db DB) Cleanup() (int64, error) {
	info := db.db.Delete(File{}, "death < ?", time.Now())
	return info.RowsAffected, info.Error
}

func (db DB) AddMeta(lifespan int, hashes []string) (string, bool) {
	death := time.Now().Add(time.Duration(lifespan) * time.Second)
	short := randSeq(10)

	files := make([]File, len(hashes))
	for index, hash := range hashes {
		db.db.First(&files[index], "hash = ?", hash)
	}

	metadata := Metadata{Files: files, Short: short, Death: death}
	db.db.Create(&metadata)

	if db.db.NewRecord(metadata) {
		return "", false
	}
	return short, true
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
