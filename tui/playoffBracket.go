package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

// Column layout (left to right):
// 0=East R1 (4 series), 1=East Semis (2), 2=East Finals (1),
// 3=Finals (1), 4=West Finals (1), 5=West Semis (2), 6=West R1 (4)
var colSeriesCount = [7]int{4, 2, 1, 1, 1, 2, 4}

// cursorIndexForColRow maps (column, row) to a series index in PlayoffBracket.Series
func cursorIndexForColRow(col, row int) int {
	switch col {
	case 0:
		return row
	case 1:
		return 4 + row
	case 2:
		return 6
	case 3:
		return 7
	case 4:
		return 8
	case 5:
		return 9 + row
	case 6:
		return 11 + row
	}
	return 0
}

// colRowForIndex is the inverse of cursorIndexForColRow
func colRowForIndex(idx int) (col, row int) {
	switch {
	case idx <= 3:
		return 0, idx
	case idx <= 5:
		return 1, idx - 4
	case idx == 6:
		return 2, 0
	case idx == 7:
		return 3, 0
	case idx == 8:
		return 4, 0
	case idx <= 10:
		return 5, idx - 9
	default:
		return 6, idx - 11
	}
}

type PlayoffBracket struct {
	seasonSelector SeasonSelector
	bracket        types.PlayoffBracket
	cursorCol      int
	cursorRow      int
	season         string
	loading        bool
	width          int
	height         int
	quitting       bool
}

type bracketFetchedMsg struct {
	bracket types.PlayoffBracket
	err     error
}

// NewPlayoffBracket creates the playoff bracket view for the given season,
// restoring cursor to restoreCursorIdx (series index from cursor table).
// Full implementation is in Task 11; this stub satisfies the compiler.
func NewPlayoffBracket(season string, restoreCursorIdx int, size tea.WindowSizeMsg) (*PlayoffBracket, tea.Cmd, error) {
	return nil, nil, nil
}

func (m PlayoffBracket) Init() tea.Cmd { return nil }

func (m PlayoffBracket) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m PlayoffBracket) View() string {
	return ""
}
