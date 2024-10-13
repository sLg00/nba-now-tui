package logger

import (
	"log"
	"os"
)

// LogToFile is a helper function which sets up a daily log file in the application's config directory.
// The function is quite lazy atm, with paths hardcoded. I'm ok with it.
func LogToFile() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return
	}
	if os.IsNotExist(err) == true {
		err = os.Mkdir(home+"/.config/nba-tui/logs/", 0777)
		if err != nil {
			log.Println(err)
			return
		}
	}
	fileName := home + "/.config/nba-tui/logs/appLog.log"
	_, err = os.Stat(fileName)

	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
		return
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Initiating logger. Logging to: ", fileName)
}
