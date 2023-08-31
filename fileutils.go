package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func deleteFilesInDirectory(dirPath string) error {

	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0755) // 0755 permission
		if err != nil {
			return err
		}
		fmt.Printf("Directory '%s' created.\n", dirPath)
	} else if err != nil {
		return err
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	if len(files) > 0 {
		for _, file := range files {
			filePath := filepath.Join(dirPath, file.Name())
			err := os.Remove(filePath)
			if err != nil {
				fmt.Printf(`Error deleting images: %v`, err)
			}
		}
	}

	return nil
}
