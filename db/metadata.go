package db

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

type Metadata struct {
	gorm.Model
	Files []File
	Short string    `gorm:"not null;unique"`
	Death time.Time `gorm:"not null"`
}

func AddMeta(lifespan int, hashes []string) (string, bool) {
	death := time.Now().Add(time.Duration(lifespan) * time.Second)
	short := randSeq(10)

	files := make([]File, len(hashes))
	for index, hash := range hashes {
		db.First(&files[index], "hash = ?", hash)
	}

	metadata := Metadata{Files: files, Short: short, Death: death}
	db.Create(&metadata)

	if db.NewRecord(metadata) {
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
