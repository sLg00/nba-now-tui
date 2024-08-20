package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/evertras/bubble-table/table"
)

type boxScore struct {
	placeholder    viewport.Model
	placeholderTwo table.Model
}

//func (m boxScore) Init() tea.Cmd { return nil }
//
//func initBoxScore(gameID string, p *tea.Program) (*boxScore, error) {
//
//	//take teh unmarshalled json
//	//convert the json objects to strings
//	//draw UI
//}
