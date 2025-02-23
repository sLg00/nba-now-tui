package filesystemops

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// FindFiles takes a path to a directory and a regexp pattern (as a string) and returns a list of matching files.
// Directories are skipped.
func FindFiles(path string, pattern string) ([]string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("path %s does not exist ", path)
	}

	dirContents, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	validatedRegex := regexp.MustCompile(pattern)

	var fileList []string
	for _, item := range dirContents {
		if !item.IsDir() {
			if validatedRegex.MatchString(item.Name()) {
				filePath := filepath.Join(path, item.Name())

				fileInfo, err := os.Stat(filePath)
				if err != nil {
					log.Println("error getting file info for ", filePath)
					continue
				}

				if time.Since(fileInfo.ModTime()) > 72*time.Hour {
					fileList = append(fileList, filePath)
				}
			}
		}
	}
	return fileList, nil
}

// RemoveFiles takes a list of valid FULL file paths and removes said files.
func RemoveFiles(fileList []string) error {
	for _, file := range fileList {
		err := os.Remove(file)
		if err != nil {
			log.Println("could not remove files")
			return err
		}
	}
	return nil
}
