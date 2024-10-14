package main

import (
	"github.com/sLg00/nba-now-tui/cmd/internal"
	"github.com/sLg00/nba-now-tui/tui"
)

func main() {
	internal.LogToFile()
	//err := client.NewClient().MakeDefaultRequests()
	//if err != nil {
	//	log.Println(err)
	//	os.Exit(1)
	//}
	tui.RenderUI()

}
