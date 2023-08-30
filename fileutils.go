package main

import (
	"os"
	"path/filepath"
)

func deleteFilesInDirectory(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}
