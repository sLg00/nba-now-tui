package client

import (
	"log"
	"os"
	"time"
)

type PathComponents struct {
	Home    string //home directory of the current OS user
	Path    string //path to the config directory of the app
	LLFile  string //league leaders file name
	SSFile  string //season standings file name
	DSBFile string //daily scoreboard file name
}

func (p PathComponents) LLFullPath() string {
	return p.Home + p.Path + p.LLFile
}

func (p PathComponents) SSFullPath() string {
	return p.Home + p.Path + p.SSFile
}

func (p PathComponents) DSBFullPath() string {
	return p.Home + p.Path + p.DSBFile
}

// InstantiatePaths is a factory function that returns a PathComponents struct with default values
func InstantiatePaths() PathComponents {
	today := time.Now().Format("2006-01-02")
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
	}
	paths := PathComponents{
		Home:    home,
		Path:    "/.config/nba-tui/",
		LLFile:  today + "_ll",
		SSFile:  today + "_ss",
		DSBFile: today + "_dsb",
	}
	return paths
}

// createDirectory creates the dir to hold daily json files received from the NBA API. If a directory already exists,
// nothing his done
func createDirectory(pc PathComponents) (string, error) {
	path := pc.Home + pc.Path

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
	paths := NewClient().InstantiatePaths()
	_, err := createDirectory(paths)
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
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}
