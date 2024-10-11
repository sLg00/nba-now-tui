package logger

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/tui"
	"log"
	"os"
)

// LogToFile is a helper function which sets up a daily log file in the application's config directory.
// The function is quite lazy atm, with paths hardcoded. I'm ok with it.
func LogToFile() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
	}
	fileName := home + "/.config/nba-tui/logs/log_" + tui.Date() + ".log"
	_, err = os.Stat(fileName)
	if os.IsNotExist(err) == true {
		err := os.Mkdir(home+"/.config/nba-tui/logs/", 0777)
		if err != nil {
			err = fmt.Errorf("Could not create file: %w\n", err)
			log.Println(err)
		}
	}
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Initiating logger. Logging to: ", fileName)
}
