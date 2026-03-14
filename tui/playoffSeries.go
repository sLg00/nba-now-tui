package tui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sLg00/nba-now-tui/cmd/converters"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

type PlayoffSeriesView struct {
	series        types.PlayoffSeries
	games         []types.PlayoffGame
	cursor        int
	bracketCursor int
	season        string
	width         int
	height        int
	loading       bool
	quitting      bool
}

type playoffSeriesGamesFetchedMsg struct {
	games []types.PlayoffGame
	err   error
}

func NewPlayoffSeries(series types.PlayoffSeries, bracketCursor int, season string, size tea.WindowSizeMsg) (*PlayoffSeriesView, tea.Cmd, error) {
	m := &PlayoffSeriesView{
		series:        series,
		bracketCursor: bracketCursor,
		season:        season,
		width:         size.Width,
		height:        size.Height,
		loading:       true,
	}
	return m, fetchPlayoffSeriesGamesCmd(season), nil
}

func fetchPlayoffSeriesGamesCmd(season string) tea.Cmd {
	return func() tea.Msg {
		cl := nbaAPI.NewClient()
		if err := cl.FetchCommonPlayoffSeries(season); err != nil {
			return playoffSeriesGamesFetchedMsg{err: err}
		}
		rs, err := cl.Loader.LoadCommonPlayoffSeries(season)
		if err != nil {
			return playoffSeriesGamesFetchedMsg{err: err}
		}
		games, err := converters.PopulatePlayoffSeriesGames(rs, "")
		return playoffSeriesGamesFetchedMsg{games: games, err: err}
	}
}

func (m PlayoffSeriesView) Init() tea.Cmd { return nil }

func (m PlayoffSeriesView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case playoffSeriesGamesFetchedMsg:
		m.loading = false
		if msg.err != nil {
			log.Println("error fetching playoff series games:", msg.err)
			return m, nil
		}
		m.games = msg.games
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		switch {
		case key.Matches(msg, Keymap.Back):
			pb, cmd, err := NewPlayoffBracket(m.season, m.bracketCursor, WindowSize)
			if err != nil {
				log.Println(err)
				return InitMenu()
			}
			return pb, cmd
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, Keymap.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, Keymap.Down):
			if m.cursor < len(m.games)-1 {
				m.cursor++
			}
		case key.Matches(msg, Keymap.Enter):
			if m.cursor < len(m.games) {
				g := m.games[m.cursor]
				if g.Completed {
					bx, cmd, err := NewBoxScore(g.GameID, g.Date, 3, "playoffBracket", m.season, m.bracketCursor, WindowSize)
					if err != nil {
						log.Println(err)
						return m, nil
					}
					return bx, cmd
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m PlayoffSeriesView) View() string {
	if m.quitting {
		return ""
	}
	if m.loading {
		return lipgloss.NewStyle().Width(m.width).Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).Render("Loading series...")
	}

	header := fmt.Sprintf("%s vs %s — %s", m.series.TopTeam.Tricode, m.series.BottomTeam.Tricode, m.season)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))

	colStyle := lipgloss.NewStyle().Width(12)
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	dimmed := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	var rows []string
	rows = append(rows, headerStyle.Render(header))
	rows = append(rows, "")

	for i, g := range m.games {
		score := "Not played"
		if g.Completed {
			score = fmt.Sprintf("%s %d — %d %s", g.AwayTricode, g.AwayScore, g.HomeScore, g.HomeTricode)
		}
		row := fmt.Sprintf("%s  %s  %s",
			colStyle.Render(fmt.Sprintf("Game %d", g.GameNumber)),
			colStyle.Render(g.Date),
			score,
		)
		if i == m.cursor {
			row = cursorStyle.Render("▶ ") + row
		} else if !g.Completed {
			row = "  " + dimmed.Render(row)
		} else {
			row = "  " + row
		}
		rows = append(rows, row)
	}

	if len(m.games) == 0 {
		rows = append(rows, dimmed.Render("No games available"))
	}

	rows = append(rows, "", HelpStyle(HelpFooter()))
	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return DocStyle.Render(content)
}
