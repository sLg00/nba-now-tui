package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/app/datamodels"
	"log"
	"strings"
)

type dailyView struct {
	gameCards  []table.Model
	quitting   bool
	focusIndex int
	numCols    int
}

func newGameCard(r []table.Row) table.Model {
	columns := []table.Column{
		table.NewColumn("teams", "", 5),
		table.NewColumn("scores", "", 5),
	}

	gc := table.New(columns)

	gc = gc.WithRows(r)
	gc = gc.BorderDefault()
	gc = gc.WithHeaderVisibility(false)
	return gc
}

func (m dailyView) Init() tea.Cmd { return nil }

func initDailyView(i list.Item, p *tea.Program) (*dailyView, error) {
	dailyScores, _, err := datamodels.PopulateDailyGameResults()
	if err != nil {
		log.Printf("Could not populate daily results: %v", err)
		return nil, err
	}
	dailyScoresStrings := datamodels.ConvertToString(dailyScores)

	var gameCards []table.Model

	for _, gameScore := range dailyScoresStrings {
		var scoreRows []table.Row

		homeRowData := table.RowData{
			"teams":  gameScore[4],
			"scores": gameScore[3],
		}
		awayRowData := table.RowData{
			"teams":  gameScore[8],
			"scores": gameScore[7],
		}

		scoreRows = append(scoreRows, table.NewRow(homeRowData))
		scoreRows = append(scoreRows, table.NewRow(awayRowData))

		gameCard := newGameCard(scoreRows)
		gameCards = append(gameCards, gameCard)
	}

	m := &dailyView{gameCards: gameCards, focusIndex: 0, numCols: 3}
	return m, nil
}

func (m dailyView) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Back):
			return InitMenu()
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, Keymap.Enter):
			//open box score
		case key.Matches(msg, Keymap.Up):
			if m.focusIndex >= m.numCols {
				m.focusIndex -= m.numCols
			}
		case key.Matches(msg, Keymap.Down):
			if m.focusIndex+m.numCols < len(m.gameCards) {
				m.focusIndex += m.numCols
			}
		case key.Matches(msg, Keymap.Left):
			if m.focusIndex%m.numCols > 0 {
				m.focusIndex--
			}
		case key.Matches(msg, Keymap.Right):
			if m.focusIndex%m.numCols < m.numCols-1 && m.focusIndex+1 < len(m.gameCards) {
				m.focusIndex++
			}
		}
	}
	m.gameCards[m.focusIndex], cmd = m.gameCards[m.focusIndex].Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m dailyView) View() string {
	if m.quitting {
		return ""
	}

	var focusedBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("5")). // Purple color code
		Padding(1)

	var b strings.Builder

	for i, gameCard := range m.gameCards {
		if i == m.focusIndex {
			// Apply purple border to the focused gameCard
			gameCard = gameCard.WithBaseStyle(focusedBorderStyle)
		}

		b.WriteString(gameCard.View())
		b.WriteString("\n")

		if (i+1)%m.numCols == 0 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
