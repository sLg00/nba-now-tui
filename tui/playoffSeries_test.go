package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

func TestPlayoffSeries_Init_IssuesCmd(t *testing.T) {
	series := types.PlayoffSeries{
		SeriesID:   "0042300401",
		TopTeam:    types.PlayoffTeam{Tricode: "BOS"},
		BottomTeam: types.PlayoffTeam{Tricode: "MIA"},
	}
	ps, cmd, err := NewPlayoffSeries(series, 0, "2023-24", tea.WindowSizeMsg{Width: 120, Height: 40})
	if err != nil {
		t.Fatalf("NewPlayoffSeries() error: %v", err)
	}
	if ps == nil {
		t.Fatal("NewPlayoffSeries() returned nil model")
	}
	if cmd == nil {
		t.Error("NewPlayoffSeries() returned nil cmd, expected fetch command")
	}
}
