package internal

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
			match := validatedRegex.FindStringSubmatch(item.Name())
			if len(match) > 0 {
				fileDateStr := match[1]
				fileDate, err := time.Parse("2006-01-02", fileDateStr)
				if err != nil {
					log.Println("error parsing file date: " + fileDateStr)
				}

				if fileDate.Before(time.Now().Add(-48 * time.Hour)) {
					filePath := filepath.Join(path, item.Name())
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
