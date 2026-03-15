package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

const tableNodeWidth = 10 // "(1)BOS 4  " — padded to 10 chars

var (
	bracketHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	bracketDimStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	bracketNormalStyle    = lipgloss.NewStyle()
)

type bracketTableRenderer struct {
	series [15]types.PlayoffSeries
	cursor int
}

func newBracketTableRenderer(series []types.PlayoffSeries, cursor int) bracketTableRenderer {
	br := bracketTableRenderer{cursor: cursor}
	for i := 0; i < 15 && i < len(series); i++ {
		br.series[i] = series[i]
	}
	return br
}

// renderTeamCell formats a single team slot, applying highlighting or dim style.
func renderTeamCell(t types.PlayoffTeam, highlighted bool) string {
	if t.Tricode == "" || t.Tricode == "TBD" {
		return bracketDimStyle.Render(fmt.Sprintf("%-*s", tableNodeWidth, "TBD"))
	}
	style := bracketNormalStyle
	if highlighted {
		style = bracketHighlightStyle
	}
	var label string
	if t.Seed > 0 {
		label = fmt.Sprintf("(%d)%-3s %d", t.Seed, t.Tricode, t.Wins)
	} else {
		label = fmt.Sprintf("%-3s %d", t.Tricode, t.Wins)
	}
	return style.Render(fmt.Sprintf("%-*s", tableNodeWidth, label))
}

// Render produces the columnar bracket table as a multi-line string.
//
// Column order (left to right): W R1 | W Semis | W CF | Finals | E CF | E Semis | E R1
//
// Row placement (8 data lines, 2 per R1 slot):
//
//	0: W_R1[0]top,                              E_R1[0]top
//	1: W_R1[0]bot, W_Semi[0]top,                E_Semi[0]top, E_R1[0]bot
//	2: W_R1[1]top,              W_CF_top,        E_CF_top,     E_R1[1]top
//	3: W_R1[1]bot, W_Semi[0]bot, Finals_top,     E_Semi[0]bot, E_R1[1]bot
//	4: W_R1[2]top,               Finals_bot,                   E_R1[2]top
//	5: W_R1[2]bot, W_Semi[1]top, W_CF_bot,       E_CF_bot,     E_Semi[1]top, E_R1[2]bot
//	6: W_R1[3]top,                                             E_R1[3]top
//	7: W_R1[3]bot, W_Semi[1]bot,                 E_Semi[1]bot, E_R1[3]bot
func (br bracketTableRenderer) Render() string {
	s := br.series
	c := br.cursor

	// Alias series by role for readability.
	// Series index layout: 0-3=E R1, 4-5=E Semis, 6=E CF, 7=Finals, 8=W CF, 9-10=W Semis, 11-14=W R1
	wR1   := [4]types.PlayoffSeries{s[11], s[12], s[13], s[14]}
	wSemi := [2]types.PlayoffSeries{s[9], s[10]}
	wCF   := s[8]
	fin   := s[7]
	eCF   := s[6]
	eSemi := [2]types.PlayoffSeries{s[4], s[5]}
	eR1   := [4]types.PlayoffSeries{s[0], s[1], s[2], s[3]}

	wR1Sel   := [4]bool{c == 11, c == 12, c == 13, c == 14}
	wSemiSel := [2]bool{c == 9, c == 10}
	wCFSel   := c == 8
	finSel   := c == 7
	eCFSel   := c == 6
	eSemiSel := [2]bool{c == 4, c == 5}
	eR1Sel   := [4]bool{c == 0, c == 1, c == 2, c == 3}

	blank := strings.Repeat(" ", tableNodeWidth)
	sep   := " │ "

	// cellAt returns the cell string for a given line and column role.
	cellAt := func(line int) (wR1c, wSc, wCFc, fc, eCFc, eSc, eR1c string) {
		wR1c, wSc, wCFc, fc, eCFc, eSc, eR1c = blank, blank, blank, blank, blank, blank, blank

		// W R1: top on even lines 0/2/4/6, bot on odd lines 1/3/5/7
		idx := line / 2
		if line%2 == 0 {
			wR1c = renderTeamCell(wR1[idx].TopTeam, wR1Sel[idx])
			eR1c = renderTeamCell(eR1[idx].TopTeam, eR1Sel[idx])
		} else {
			wR1c = renderTeamCell(wR1[idx].BottomTeam, wR1Sel[idx])
			eR1c = renderTeamCell(eR1[idx].BottomTeam, eR1Sel[idx])
		}

		switch line {
		case 1:
			wSc = renderTeamCell(wSemi[0].TopTeam, wSemiSel[0])
			eSc = renderTeamCell(eSemi[0].TopTeam, eSemiSel[0])
		case 2:
			wCFc = renderTeamCell(wCF.TopTeam, wCFSel)
			eCFc = renderTeamCell(eCF.TopTeam, eCFSel)
		case 3:
			wSc = renderTeamCell(wSemi[0].BottomTeam, wSemiSel[0])
			fc  = renderTeamCell(fin.TopTeam, finSel)
			eSc = renderTeamCell(eSemi[0].BottomTeam, eSemiSel[0])
		case 4:
			fc = renderTeamCell(fin.BottomTeam, finSel)
		case 5:
			wSc  = renderTeamCell(wSemi[1].TopTeam, wSemiSel[1])
			wCFc = renderTeamCell(wCF.BottomTeam, wCFSel)
			eCFc = renderTeamCell(eCF.BottomTeam, eCFSel)
			eSc  = renderTeamCell(eSemi[1].TopTeam, eSemiSel[1])
		case 7:
			wSc = renderTeamCell(wSemi[1].BottomTeam, wSemiSel[1])
			eSc = renderTeamCell(eSemi[1].BottomTeam, eSemiSel[1])
		}
		return
	}

	var sb strings.Builder
	sb.WriteString(buildTableHeader())
	sb.WriteByte('\n')
	for line := 0; line < 8; line++ {
		wR1c, wSc, wCFc, fc, eCFc, eSc, eR1c := cellAt(line)
		sb.WriteString(wR1c + " " + wSc + " " + wCFc + sep + fc + sep + eCFc + " " + eSc + " " + eR1c)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildTableHeader() string {
	n := tableNodeWidth
	blank := strings.Repeat(" ", n)
	sep   := " │ "

	// Label widths must match column widths used in data rows.
	// Each column is tableNodeWidth chars. Columns: W_R1, W_Semis, W_CF | Finals | E_CF, E_Semis, E_R1
	westLabel := fmt.Sprintf("%-*s", n*3+2, "WEST")
	eastLabel := fmt.Sprintf("%-*s", n*3+2, "EAST")
	finLabel  := fmt.Sprintf("%-*s", n, "FINALS")

	header1 := westLabel + sep + finLabel + sep + eastLabel
	header2 := fmt.Sprintf("%-*s", n, "R1") + " " +
		fmt.Sprintf("%-*s", n, "Semis") + " " +
		fmt.Sprintf("%-*s", n, "CF") +
		sep + blank + sep +
		fmt.Sprintf("%-*s", n, "CF") + " " +
		fmt.Sprintf("%-*s", n, "Semis") + " " +
		fmt.Sprintf("%-*s", n, "R1")
	separator := strings.Repeat("─", len([]rune(header2)))

	return header1 + "\n" + header2 + "\n" + separator
}
