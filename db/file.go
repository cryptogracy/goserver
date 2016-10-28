package db

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type File struct {
	gorm.Model
	Hash  string    `gorm:"primary_key;unique"`
	Death time.Time `gorm:"not null"`
	Short string
}

func AddFile(hash string, lifespan int) bool {
	death := time.Now().Add(time.Duration(lifespan) * time.Second)
	short := ""

	file := File{Hash: hash, Death: death, Short: short}
	db.Create(&file)

	if db.NewRecord(file) {
		return false
	}
	return true
}
