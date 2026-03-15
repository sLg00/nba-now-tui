package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sLg00/nba-now-tui/cmd/converters"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

// Column layout (left to right):
// 0=West R1 (4 series), 1=West Semis (2), 2=West Finals (1),
// 3=Finals (1), 4=East Finals (1), 5=East Semis (2), 6=East R1 (4)
var colSeriesCount = [7]int{4, 2, 1, 1, 1, 2, 4}

// cursorIndexForColRow maps (column, row) to a series index in PlayoffBracket.Series
func cursorIndexForColRow(col, row int) int {
	switch col {
	case 0:
		return 11 + row // West R1
	case 1:
		return 9 + row // West Semis
	case 2:
		return 8 // West Finals
	case 3:
		return 7 // NBA Finals
	case 4:
		return 6 // East Finals
	case 5:
		return 4 + row // East Semis
	case 6:
		return row // East R1
	}
	return 0
}

// colRowForIndex is the inverse of cursorIndexForColRow
func colRowForIndex(idx int) (col, row int) {
	switch {
	case idx <= 3: // East R1
		return 6, idx
	case idx <= 5: // East Semis
		return 5, idx - 4
	case idx == 6: // East Finals
		return 4, 0
	case idx == 7: // NBA Finals
		return 3, 0
	case idx == 8: // West Finals
		return 2, 0
	case idx <= 10: // West Semis
		return 1, idx - 9
	default: // West R1
		return 0, idx - 11
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
// restoring cursor to restoreCursorIdx (series index per cursor table).
func NewPlayoffBracket(season string, restoreCursorIdx int, size tea.WindowSizeMsg) (*PlayoffBracket, tea.Cmd, error) {
	currentSeason := nbaAPI.NewClient().Dates.GetCurrentSeason()
	ss := NewSeasonSelector(currentSeason)
	ss.season = season
	ss.SetWidth(size.Width)

	m := &PlayoffBracket{
		seasonSelector: ss,
		season:         season,
		loading:        true,
		width:          size.Width,
		height:         size.Height,
	}
	m.cursorCol, m.cursorRow = colRowForIndex(restoreCursorIdx)

	return m, fetchPlayoffBracketCmd(season), nil
}

func fetchPlayoffBracketCmd(season string) tea.Cmd {
	return func() tea.Msg {
		cl := nbaAPI.NewClient()

		if err := cl.FetchCommonPlayoffSeries(season); err != nil {
			log.Printf("fetchPlayoffBracket: FetchCommonPlayoffSeries(%s) failed: %v", season, err)
		} else if rs, err := cl.Loader.LoadCommonPlayoffSeries(season); err != nil {
			log.Printf("fetchPlayoffBracket: LoadCommonPlayoffSeries(%s) failed: %v", season, err)
		} else if bracket, err := converters.PopulatePlayoffBracket(rs, season); err != nil {
			log.Printf("fetchPlayoffBracket: PopulatePlayoffBracket(%s) failed: %v", season, err)
		} else {
			return bracketFetchedMsg{bracket: bracket}
		}

		// Fallback: project bracket from current standings.
		log.Printf("fetchPlayoffBracket: using standings projection for %s", season)
		rs, err := cl.Loader.LoadSeasonStandings()
		if err != nil {
			return bracketFetchedMsg{err: err}
		}
		teams, _, err := converters.PopulateTeamStats(rs)
		if err != nil {
			return bracketFetchedMsg{err: err}
		}
		east, west := teams.SplitStandingsPerConference()
		return bracketFetchedMsg{bracket: converters.ProjectedBracketFromStandings(east, west, season)}
	}
}

func (m PlayoffBracket) Init() tea.Cmd { return nil }

func (m PlayoffBracket) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case bracketFetchedMsg:
		m.loading = false
		if msg.err != nil {
			log.Println("error fetching playoff bracket:", msg.err)
			return m, nil
		}
		m.bracket = msg.bracket
		return m, nil

	case seasonChangedMsg:
		m.season = msg.season
		m.loading = true
		m.cursorCol, m.cursorRow = 0, 0
		return m, fetchPlayoffBracketCmd(msg.season)

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		switch {
		case key.Matches(msg, Keymap.Back):
			return InitMenu()
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, Keymap.Tab):
			if m.seasonSelector.focused {
				m.seasonSelector.Blur()
			} else {
				m.seasonSelector.Focus()
			}
		case key.Matches(msg, Keymap.Left):
			if m.seasonSelector.focused {
				var ssCmd tea.Cmd
				m.seasonSelector, ssCmd = m.seasonSelector.Update(msg)
				return m, ssCmd
			}
			if m.cursorCol > 0 {
				m.cursorCol--
				maxRow := colSeriesCount[m.cursorCol] - 1
				if m.cursorRow > maxRow {
					m.cursorRow = maxRow
				}
			}
		case key.Matches(msg, Keymap.Right):
			if m.seasonSelector.focused {
				var ssCmd tea.Cmd
				m.seasonSelector, ssCmd = m.seasonSelector.Update(msg)
				return m, ssCmd
			}
			if m.cursorCol < 6 {
				m.cursorCol++
				maxRow := colSeriesCount[m.cursorCol] - 1
				if m.cursorRow > maxRow {
					m.cursorRow = maxRow
				}
			}
		case key.Matches(msg, Keymap.Up):
			if m.cursorRow > 0 {
				m.cursorRow--
			}
		case key.Matches(msg, Keymap.Down):
			if m.cursorRow < colSeriesCount[m.cursorCol]-1 {
				m.cursorRow++
			}
		case key.Matches(msg, Keymap.Enter):
			idx := cursorIndexForColRow(m.cursorCol, m.cursorRow)
			if idx < len(m.bracket.Series) {
				series := m.bracket.Series[idx]
				if series.Status != "pre" {
					ps, cmd, err := NewPlayoffSeries(series, idx, m.season, WindowSize)
					if err != nil {
						log.Println(err)
						return m, nil
					}
					return ps, cmd
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.seasonSelector.SetWidth(msg.Width)
	}
	return m, nil
}

func (m PlayoffBracket) View() string {
	if m.quitting {
		return ""
	}
	if m.loading {
		content := m.seasonSelector.View() + "\n\nLoading bracket..."
		return lipgloss.NewStyle().Width(m.width).Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).Render(content)
	}

	cursorIdx := cursorIndexForColRow(m.cursorCol, m.cursorRow)
	br := newBracketTableRenderer(m.bracket.Series, cursorIdx)
	bracket := br.Render()

	content := lipgloss.JoinVertical(lipgloss.Left,
		m.seasonSelector.View(),
		"",
		bracket,
		"",
		HelpStyle(HelpFooter()),
	)
	return DocStyle.Render(content)
}
