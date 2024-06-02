package client

import (
	"log"
	"os"
)

var (
	//OS      = runtime.GOOS
	HOME, _ = os.UserHomeDir()
	PATH    = "/.config/nba-tui/"
	LLFILE  = Today + "_ll"
	SSFILE  = Today + "_ss"
	DSBFILE = Today + "_dsb"

	LLFULLPATH  = HOME + PATH + LLFILE
	SSFULLPATH  = HOME + PATH + SSFILE
	DSBFULLPATH = HOME + PATH + DSBFILE
)

// createDirectory creates the dir to hold daily json files received from the NBA API. If a directory already exists,
// nothing his done
func createDirectory() (string, error) {
	path := HOME + PATH

	_, err := os.Stat(path)
	if os.IsNotExist(err) == true {

		err = os.Mkdir(path, 0777)
		if err != nil {
			log.Println("Error creating directory", err)
			return path, err
		}
	}
	return path, nil
}

// WriteToFiles handles the writing of the json responses to the filesystem. It takes a string
// (the full path of the file the body of the JSON response as []byte
func WriteToFiles(s string, b []byte) error {
	_, err := createDirectory()
	if err != nil {
		log.Println(err)
	}
	filePath := s

	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(b)
	if err != nil {
		log.Println("could not write to file")
		return err
	}
	return nil
}

// fileChecker is a helper function to check if a file exists in the system
func fileChecker(s string) bool {
	_, err := os.Stat(s)
	if os.IsExist(err) == true {
		return true
	} else {
		return false
	}
}
