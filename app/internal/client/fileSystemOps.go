package client

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	HOME, _ = os.UserHomeDir()
)

// createDirectory creates the dir to hold daily json files received from the NBA API. If a directory already exists,
// nothing his done
func createDirectory() string {
	path := HOME + "/.config/nba-tui/"

	_, err := os.Stat(path)
	if os.IsNotExist(err) == true {

		err = os.Mkdir(path, 0777)
		if err != nil {
			fmt.Println("Error creating directory", err)
		}
	}
	return path
}

// WriteToFiles handles the writing of the json responses to the filesystem
func WriteToFiles() error {
	path := createDirectory()
	llFile := "ll_" + Today
	filePath := filepath.Join(path, llFile)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(LLJson)
	if err != nil {
		log.Println("could not write to file")
		return err
	}
	return nil
}
