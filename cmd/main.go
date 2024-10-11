package main

import (
	"github.com/sLg00/nba-now-tui/cmd/client"
	"github.com/sLg00/nba-now-tui/cmd/internal/logger"
	"github.com/sLg00/nba-now-tui/tui"
)

func main() {
	logger.LogToFile()
	client.NewClient().MakeDefaultRequests()
	tui.RenderUI()

}
