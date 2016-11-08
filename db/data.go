package db

import (
	"log"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func Init(database string) error {

	rand.Seed(time.Now().UTC().UnixNano())
	var err error
	db, err = gorm.Open("sqlite3", database)
	if err != nil {
		return err
	}

	db.LogMode(false)

	db.AutoMigrate(&File{}, &Metadata{})

	return nil
}

func Close() error {
	return db.Close()
}

func Cleanup() {
	if affected, err := removeFiles(); err == nil {
		log.Printf("Remove %v old files\n", affected)
	} else {
		log.Printf("Unable to delete old files", err)
	}
}
