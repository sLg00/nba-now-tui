package client

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/internal"
	"log"
	"os"
)

type PathComponents struct {
	Home            string //home directory of the current OS user
	Path            string //path to the config directory of the cmd
	LLFile          string //league leaders file name
	SSFile          string //season standings file name
	DSBFile         string //daily scoreboard file name
	BoxScorePath    string //folder to store box scores
	BoxScoreFile    string //box score file name
	TeamProfilePath string //folder to store profile pages
	TeamProfileFile string // team profile file
}

func (p PathComponents) LLFullPath() string {
	return p.Home + p.Path + p.LLFile
}

func (p PathComponents) SSFullPath() string { return p.Home + p.Path + p.SSFile }

func (p PathComponents) DSBFullPath() string {
	return p.Home + p.Path + p.DSBFile
}

func (p PathComponents) BoxScoreFullPath() string {
	return p.Home + p.Path + p.BoxScorePath + p.BoxScoreFile
}

func (p PathComponents) TeamProfileFullPath() string {
	return p.Home + p.Path + p.TeamProfilePath + p.TeamProfileFile
}

// InstantiatePaths is a factory function that returns a PathComponents struct with default values
func InstantiatePaths(s string) *PathComponents {
	today, err := GetDateArg()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("could not determine home directory: %w", err)
		log.Println(err)
	}
	paths := PathComponents{
		Home:            home,
		Path:            "/.config/nba-tui/",
		LLFile:          today + "_ll",
		SSFile:          today + "_ss",
		DSBFile:         today + "_dsb",
		BoxScorePath:    "boxscores/",
		BoxScoreFile:    today + "_" + s,
		TeamProfilePath: "teamprofiles/",
		TeamProfileFile: s,
	}
	return &paths
}

// createDirectory creates the dir to hold daily json files received from the NBA API. If a directory already exists,
// nothing is done
func createDirectory(pc *PathComponents) (string, error) {
	path := pc.Home + pc.Path + pc.BoxScorePath
	teamProfilePath := pc.Home + pc.Path + pc.TeamProfilePath
	paths := []string{path, teamProfilePath}

	for _, path = range paths {
		_, err := os.Stat(path)
		if os.IsNotExist(err) == true {
			err = os.Mkdir(path, 0777)
			if err != nil {
				err = fmt.Errorf("could not create directory: %w\n", err)
				log.Println(err)
				return path, err
			}
		}
	}
	return path, nil
}

// WriteToFiles handles the writing of the json responses to the filesystem. It takes a string
// (the full path of the file the body of the JSON response as []byte
func WriteToFiles(s string, b []byte) error {
	paths := NewClient().InstantiatePaths(s)
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
	fileInfo, err := os.Stat(s)
	if err == nil {
		if fileInfo.Size() > 1000 {
			return true
		}
	}
	return false
}

func CleanOldFiles(pc *PathComponents) error {
	dailyFilesPath := pc.Home + pc.Path
	boxScoreFilesPath := pc.Home + pc.Path + pc.BoxScorePath
	paths := []string{dailyFilesPath, boxScoreFilesPath}

	filesRegex := "^(\\d{4}-\\d{2}-\\d{2})_.*$"

	for _, p := range paths {
		fileList, err := internal.FindFiles(p, filesRegex)
		if err != nil {
			log.Println(err)
			return err
		}
		err = internal.RemoveFiles(fileList)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
