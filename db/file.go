package db

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/cryptogracy/goserver/configuration"
)

var (
	ErrHash      = errors.New("Wrong Hash")
	ErrInternal  = errors.New("Internal Error")
	ErrFileExist = errors.New("File exists")
)

type File struct {
	gorm.Model
	Hash   string    `gorm:"primary_key;unique"`
	Death  time.Time `gorm:"not null"`
	Short  string
	Reader io.Reader `gorm:"-"`
}

func AddFile(hash string, lifespan int, reader io.Reader) (err error) {

	death := time.Now().UTC().Add(time.Duration(lifespan) * time.Minute)
	short := ""

	file := File{Hash: hash, Death: death, Short: short, Reader: reader}
	err = db.Create(&file).Error
	// This is not nice, but I have no idea how to make it better
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
		err = ErrFileExist
	}
	return
}

func removeFiles() (int64, error) {
	info := db.Delete(File{}, "death < ?", time.Now().UTC())
	return info.RowsAffected, info.Error
}

func (file *File) AfterDelete() error {
	return os.Remove(path.Join(configuration.Config.Dir, file.Hash))
}

func (file *File) AfterCreate() error {
	// Open tmp file
	tempfile := path.Join(configuration.Config.Tempdir, file.Hash)
	out, err := os.OpenFile(tempfile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		log.Println("Unable to open tempfile", tempfile, err)
		return ErrInternal
	}
	defer out.Close()
	defer os.Remove(tempfile) // We never want to keep the tmpfile

	// Copy content to tmp
	if _, err := io.Copy(out, file.Reader); err != nil {
		return ErrInternal
	}

	// Check if the hash is correct
	if _, err := out.Seek(0, os.SEEK_SET); err != nil {
		return ErrInternal
	}

	if !isHash(file.Hash, out) {
		return ErrHash
	}

	// Put into configuration Dir
	if err := os.Rename(tempfile, path.Join(configuration.Config.Dir, file.Hash)); err != nil {
		return ErrInternal
	}

	return nil
}

func isHash(hash string, file io.Reader) bool {
	hasher := sha512.New()
	_, err := io.Copy(hasher, file)
	if err != nil {
		panic(err)
	}
	if fmt.Sprintf("%x", hasher.Sum(nil)) == hash {
		return true
	}
	return false
}
