package internal

import (
	"log"
	"os"
	"path/filepath"
)

// LogToFile is a helper function which sets up a log file in the application's config directory.
// The function is quite lazy atm, with paths hardcoded. I'm ok with it.
func LogToFile() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return
	}
	logDir := filepath.Join(home, ".config/nba-tui/logs/")
	if _, err = os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0644)
		if err != nil {
			log.Println(err)
			return
		}
	}
	fileName := filepath.Join(logDir, "appLog.log")
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
		return
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Initiating logger. Logging to: ", fileName)
}
