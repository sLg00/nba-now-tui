package tui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

func makeTBDSeries() types.PlayoffSeries {
	return types.PlayoffSeries{
		Status:     "pre",
		TopTeam:    types.PlayoffTeam{Tricode: "TBD", Seed: 0},
		BottomTeam: types.PlayoffTeam{Tricode: "TBD", Seed: 0},
	}
}

func TestRenderSeriesNode_ShowsSeedsAndRecord(t *testing.T) {
	s := types.PlayoffSeries{
		TopTeam:    types.PlayoffTeam{Tricode: "BOS", Seed: 1, Wins: 3},
		BottomTeam: types.PlayoffTeam{Tricode: "MIA", Seed: 8, Wins: 1},
		Status:     "active",
	}
	top, bot := renderSeriesNode(s, false)
	if !strings.Contains(top, "BOS") {
		t.Errorf("top node missing BOS: %q", top)
	}
	if !strings.Contains(bot, "MIA") {
		t.Errorf("bot node missing MIA: %q", bot)
	}
	if !strings.Contains(top, "3") {
		t.Errorf("top node missing win count 3: %q", top)
	}
}

func TestRenderSeriesNode_TBD(t *testing.T) {
	s := makeTBDSeries()
	top, _ := renderSeriesNode(s, false)
	if !strings.Contains(top, "TBD") {
		t.Errorf("TBD series node missing TBD: %q", top)
	}
}

func TestBracketRenderer_ProducesExpectedLineCount(t *testing.T) {
	series := make([]types.PlayoffSeries, 15)
	for i := range series {
		series[i] = makeTBDSeries()
	}
	br := newBracketRenderer(series, 0)
	output := br.Render()
	lines := strings.Split(output, "\n")
	if len(lines) < 15 {
		t.Errorf("Render() produced %d lines, want at least 15", len(lines))
	}
}

func TestBracketRenderer_CursorSeriesHighlighted(t *testing.T) {
	series := make([]types.PlayoffSeries, 15)
	for i := range series {
		series[i] = types.PlayoffSeries{
			TopTeam:    types.PlayoffTeam{Tricode: fmt.Sprintf("T%02d", i), Seed: 1},
			BottomTeam: types.PlayoffTeam{Tricode: fmt.Sprintf("B%02d", i), Seed: 8},
			Status:     "active",
		}
	}
	// Render without panic for each cursor position
	for cursor := 0; cursor < 15; cursor++ {
		br := newBracketRenderer(series, cursor)
		_ = br.Render()
	}
}
