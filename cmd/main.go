package main

import (
	"github.com/sLg00/nba-now-tui/cmd/internal"
	"github.com/sLg00/nba-now-tui/tui"
)

func main() {
	internal.LogToFile()
	tui.RenderUI()

}
