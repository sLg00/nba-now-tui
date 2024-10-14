package main

import (
	"github.com/sLg00/nba-now-tui/cmd/client"
	"github.com/sLg00/nba-now-tui/cmd/logger"
	"github.com/sLg00/nba-now-tui/tui"
	"log"
	"os"
)

func main() {
	logger.LogToFile()
	err := client.NewClient().MakeDefaultRequests()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	tui.RenderUI()

}
