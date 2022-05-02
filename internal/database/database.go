package database

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/0x00-ketsu/taskcli/internal/utils"
	"github.com/asdine/storm/v3"
	"github.com/mitchellh/go-homedir"
)

// Connect builds an connection to boltDB
func Connect(filePath string) *storm.DB {
	dbPath, _ := homedir.Expand(filePath)
	baseDir := filepath.Dir(dbPath)
	os.MkdirAll(baseDir, 0775)

	if err := utils.CreateFileIfNotExist(dbPath); err != nil {
		f, _ := ioutil.TempFile("taskcli", "bolt.db")
		dbPath = f.Name()
	}

	db, openErr := storm.Open(dbPath)
	if openErr != nil {
		fmt.Println("Error: open boltDB file failed")
		os.Exit(0)
	}

	return db
}
