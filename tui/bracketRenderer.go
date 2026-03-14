package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

const nodeWidth = 10 // "(1) BOS 3" = 9 chars, padded to 10

var (
	bracketHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	bracketDimStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	bracketNormalStyle    = lipgloss.NewStyle()
)

type bracketRenderer struct {
	series [15]types.PlayoffSeries
	cursor int
}

func newBracketRenderer(series []types.PlayoffSeries, cursor int) bracketRenderer {
	br := bracketRenderer{cursor: cursor}
	for i := 0; i < 15 && i < len(series); i++ {
		br.series[i] = series[i]
	}
	return br
}

// renderSeriesNode returns the top-team and bottom-team display strings for a series node.
func renderSeriesNode(s types.PlayoffSeries, highlighted bool) (top, bot string) {
	style := bracketNormalStyle
	if highlighted {
		style = bracketHighlightStyle
	}
	if s.Status == "pre" && (s.TopTeam.Tricode == "TBD" || s.TopTeam.Tricode == "") {
		style = bracketDimStyle
	}

	formatTeam := func(t types.PlayoffTeam) string {
		if t.Tricode == "TBD" || t.Tricode == "" {
			return fmt.Sprintf("%-*s", nodeWidth, "TBD")
		}
		if t.Seed > 0 {
			return fmt.Sprintf("(%d)%-3s %d", t.Seed, t.Tricode, t.Wins)
		}
		return fmt.Sprintf("%-3s %d", t.Tricode, t.Wins)
	}

	top = style.Render(fmt.Sprintf("%-*s", nodeWidth, formatTeam(s.TopTeam)))
	bot = style.Render(fmt.Sprintf("%-*s", nodeWidth, formatTeam(s.BottomTeam)))
	return
}

// Render produces the full ASCII bracket as a multi-line string.
// Layout: East (R1→Semis→CF) | Finals | West (CF→Semis→R1)
func (br bracketRenderer) Render() string {
	s := br.series
	cursor := br.cursor

	east := br.buildEastHalf(s[0:7], cursor)
	west := br.buildWestHalf(s[8:15], cursor)
	finalsTop, finalsBot := renderSeriesNode(s[7], cursor == 7)

	pad := strings.Repeat(" ", nodeWidth)
	var lines []string
	for i := 0; i < len(east); i++ {
		finalsCol := pad
		if i == 6 {
			finalsCol = finalsTop
		} else if i == 8 {
			finalsCol = finalsBot
		}
		lines = append(lines, east[i]+" "+finalsCol+" "+west[i])
	}
	return strings.Join(lines, "\n")
}

// buildEastHalf constructs East conference bracket lines (15 rows).
// series: [E_R1_1v8, E_R1_4v5, E_R1_3v6, E_R1_2v7, E_Semi_upper, E_Semi_lower, E_Finals]
func (br bracketRenderer) buildEastHalf(series []types.PlayoffSeries, cursor int) []string {
	blank := strings.Repeat(" ", nodeWidth)
	lines := make([]string, 15)
	for i := range lines {
		lines[i] = blank
	}

	if len(series) < 7 {
		return lines
	}

	r1Top0, r1Bot0 := renderSeriesNode(series[0], cursor == 0)
	r1Top1, r1Bot1 := renderSeriesNode(series[1], cursor == 1)
	r1Top2, r1Bot2 := renderSeriesNode(series[2], cursor == 2)
	r1Top3, r1Bot3 := renderSeriesNode(series[3], cursor == 3)
	sTop0, sBot0 := renderSeriesNode(series[4], cursor == 4)
	sTop1, sBot1 := renderSeriesNode(series[5], cursor == 5)
	cfTop, cfBot := renderSeriesNode(series[6], cursor == 6)

	con := " ─┐"
	mid := strings.Repeat(" ", nodeWidth) + "   ├─ "
	cfCon := strings.Repeat(" ", nodeWidth+6+nodeWidth) + "   ├─ "

	lines[0] = r1Top0 + con
	lines[1] = mid + sTop0 + con
	lines[2] = r1Bot0 + con
	lines[3] = cfCon + cfTop + con
	lines[4] = r1Top1 + con
	lines[5] = mid + sBot0 + " ─┘"
	lines[6] = r1Bot1 + " ─┘"
	lines[7] = strings.Repeat(" ", nodeWidth+6+nodeWidth+6+nodeWidth) + "   ├─"
	lines[8] = r1Top2 + con
	lines[9] = mid + sTop1 + con
	lines[10] = r1Bot2 + con
	lines[11] = cfCon + cfBot + " ─┘"
	lines[12] = r1Top3 + con
	lines[13] = mid + sBot1 + " ─┘"
	lines[14] = r1Bot3 + " ─┘"

	return lines
}

// buildWestHalf constructs West conference bracket lines (15 rows, mirrored).
// series: [W_Finals, W_Semi_upper, W_Semi_lower, W_R1_1v8, W_R1_4v5, W_R1_3v6, W_R1_2v7]
func (br bracketRenderer) buildWestHalf(series []types.PlayoffSeries, cursor int) []string {
	blank := strings.Repeat(" ", nodeWidth)
	lines := make([]string, 15)
	for i := range lines {
		lines[i] = blank
	}

	if len(series) < 7 {
		return lines
	}

	cfTop, cfBot := renderSeriesNode(series[0], cursor == 8)
	sTop0, sBot0 := renderSeriesNode(series[1], cursor == 9)
	sTop1, sBot1 := renderSeriesNode(series[2], cursor == 10)
	r1Top0, r1Bot0 := renderSeriesNode(series[3], cursor == 11)
	r1Top1, r1Bot1 := renderSeriesNode(series[4], cursor == 12)
	r1Top2, r1Bot2 := renderSeriesNode(series[5], cursor == 13)
	r1Top3, r1Bot3 := renderSeriesNode(series[6], cursor == 14)

	lines[0] = "┌─ " + r1Top0
	lines[1] = "│  " + r1Bot0
	lines[2] = "├─ " + sTop0 + " ─┘"
	lines[3] = "│   ┌─ " + cfTop
	lines[4] = "│   │  " + r1Top1
	lines[5] = "│   ├─ " + sBot0
	lines[6] = "│   │  " + r1Bot1
	lines[7] = "├───┘"
	lines[8] = "│   ┌─ " + r1Top2
	lines[9] = "│   ├─ " + sTop1
	lines[10] = "│   │  " + r1Bot2
	lines[11] = "│   └─ " + cfBot
	lines[12] = "│      " + r1Top3
	lines[13] = "└─ " + sBot1
	lines[14] = "   " + r1Bot3

	return lines
}
