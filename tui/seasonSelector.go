package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const minSeasonYear = 2000

type seasonChangedMsg struct {
	season string
}

type SeasonSelector struct {
	season  string
	ceiling string
	focused bool
	width   int
}

func NewSeasonSelector(currentSeason string) SeasonSelector {
	return SeasonSelector{
		season:  currentSeason,
		ceiling: currentSeason,
		focused: false,
	}
}

func (ss *SeasonSelector) SetWidth(w int) { ss.width = w }
func (ss *SeasonSelector) Focus()         { ss.focused = true }
func (ss *SeasonSelector) Blur()          { ss.focused = false }

func formatSeasonFromYear(startYear int) string {
	return fmt.Sprintf("%d-%02d", startYear, (startYear+1)%100)
}

func parseSeasonYear(season string) int {
	parts := strings.Split(season, "-")
	if len(parts) == 0 {
		return 0
	}
	year, _ := strconv.Atoi(parts[0])
	return year
}

func (ss *SeasonSelector) prevSeason() {
	year := parseSeasonYear(ss.season)
	if year > minSeasonYear {
		ss.season = formatSeasonFromYear(year - 1)
	}
}

func (ss *SeasonSelector) nextSeason() {
	year := parseSeasonYear(ss.season)
	ceilingYear := parseSeasonYear(ss.ceiling)
	if year < ceilingYear {
		ss.season = formatSeasonFromYear(year + 1)
	}
}

func (ss SeasonSelector) Update(msg tea.Msg) (SeasonSelector, tea.Cmd) {
	if !ss.focused {
		return ss, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Left):
			old := ss.season
			ss.prevSeason()
			if ss.season != old {
				return ss, func() tea.Msg { return seasonChangedMsg{season: ss.season} }
			}
		case key.Matches(msg, Keymap.Right):
			old := ss.season
			ss.nextSeason()
			if ss.season != old {
				return ss, func() tea.Msg { return seasonChangedMsg{season: ss.season} }
			}
		}
	}
	return ss, nil
}

func (ss SeasonSelector) View() string {
	seasonStyle := lipgloss.NewStyle().Bold(true)
	arrowStyle := lipgloss.NewStyle()
	if ss.focused {
		seasonStyle = seasonStyle.Foreground(lipgloss.Color("5"))
		arrowStyle = arrowStyle.Foreground(lipgloss.Color("5"))
	} else {
		arrowStyle = arrowStyle.Foreground(lipgloss.Color("241"))
	}
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("  [tab]")
	display := arrowStyle.Render("◀") + "  " + seasonStyle.Render(ss.season) + "  " + arrowStyle.Render("▶")
	if !ss.focused {
		display += hint
	}
	return lipgloss.NewStyle().Width(ss.width).Align(lipgloss.Center).Render(display)
}
