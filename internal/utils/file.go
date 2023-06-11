package utils

import (
	"io/ioutil"
	"os"
)

// CreateFileIfNotExist creates a file (absoulte path)
// returns error when failed
func CreateFileIfNotExist(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		if _, err = os.Create(filePath); err != nil {
			return err
		}
	}
	return err
}

// WriteToTempFile writes given content to a tmpFile and returns the filename
func WriteToTempFile(content string, tempFilePattern string) (string, error) {
	tmpFile, err := ioutil.TempFile("", tempFilePattern)
	if err != nil {
		return "", err
	}

	fileName := tmpFile.Name()
	if err = os.WriteFile(fileName, []byte(content), 0777); err != nil {
		return "", err
	}
	return fileName, tmpFile.Close()
}
