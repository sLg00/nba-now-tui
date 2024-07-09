package main

import (
	"github.com/sLg00/nba-now-tui/app/internal/client"
	"github.com/sLg00/nba-now-tui/tui"
)

func main() {

	client.NewClient().MakeRequests()
	tui.RenderUI()

}
