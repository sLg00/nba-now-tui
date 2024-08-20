package main

import (
	"github.com/sLg00/nba-now-tui/cmd/client"
	"github.com/sLg00/nba-now-tui/tui"
)

func main() {
	client.NewClient().MakeDefaultRequests()
	tui.RenderUI()

}
