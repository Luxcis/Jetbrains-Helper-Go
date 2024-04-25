package helper

import (
	"encoding/json"
	"log"
	"os"
)

func ReadJson(path string, payload interface{}) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return err
	}
	return nil
}

func OpenFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Failed to open or create file: %v", err)
	}
	return file
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
