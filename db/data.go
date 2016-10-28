package db

import (
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

func Cleanup() (int64, error) {
	info := db.Delete(File{}, "death < ?", time.Now())
	return info.RowsAffected, info.Error
}
