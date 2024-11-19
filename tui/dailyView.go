package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/cmd/client"
	"github.com/sLg00/nba-now-tui/cmd/datamodels"
	"log"
	"os"
)

type dailyView struct {
	gameCards  []table.Model
	quitting   bool
	focusIndex int
	numCols    int
	width      int
	height     int
}

func newGameCard(r []table.Row) table.Model {
	columns := []table.Column{
		table.NewColumn("teams", "", 7),
		table.NewColumn("scores", "", 7),
	}

	gc := table.New(columns)
	gc = gc.WithRows(r).WithBaseStyle(lipgloss.NewStyle().
		BorderStyle(lipgloss.HiddenBorder()).
		BorderForeground(lipgloss.Color("7"))).
		WithHeaderVisibility(false)
	gc = gc.BorderRounded()
	return gc
}

func (m dailyView) Init() tea.Cmd { return nil }

func initDailyView() (*dailyView, error) {
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
			"gameID": gameScore[0],
		}
		awayRowData := table.RowData{
			"teams":  gameScore[8],
			"scores": gameScore[7],
			"gameID": gameScore[0],
		}

		scoreRows = append(scoreRows, table.NewRow(homeRowData))
		scoreRows = append(scoreRows, table.NewRow(awayRowData))

		gameCard := newGameCard(scoreRows)
		gameCards = append(gameCards, gameCard)

		// For each gameID, query the NBA API, get the box score and save it to the filesystem
		err = client.NewClient().MakeOnDemandRequests(gameScore[0])
		if err != nil {
			log.Printf("Could not make on-demand requests: %v", err)
		}
	}

	m := &dailyView{gameCards: gameCards, focusIndex: 0, numCols: 3}
	return m, nil
}

func (m dailyView) getGameId() string {
	focusedCard := m.gameCards[m.focusIndex]
	rows := focusedCard.GetVisibleRows()
	if len(rows) == 1 {
		gameId, ok := rows[0].Data["gameID"].(string)
		if ok {
			return gameId
		}
	}
	if len(rows) > 1 || len(rows) < 1 {
		log.Println("Either 0 rows or more than 1 row were selected")
		//TODO: Display pop-up with User error! :)
	}
	return ""
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
			gameID := m.getGameId()
			bx, err := initBoxScore(gameID, Program)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			return bx.Update(WindowSize)
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
		case key.Matches(msg, Keymap.Tab):
			m.focusIndex = (m.focusIndex + 1) % len(m.gameCards)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.gameCards[m.focusIndex], cmd = m.gameCards[m.focusIndex].Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m dailyView) helpView() string {

	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m dailyView) View() string {
	if m.quitting {
		return ""
	}

	var rows []string
	var currentRow []string

	for i, gameCard := range m.gameCards {
		if i == m.focusIndex {
			gameCard = gameCard.WithBaseStyle(lipgloss.NewStyle().
				BorderStyle(lipgloss.HiddenBorder()).
				BorderForeground(lipgloss.Color("5"))).
				WithHeaderVisibility(false)
		}

		currentRow = append(currentRow, gameCard.View())

		if (i+1)%m.numCols == 0 {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
			currentRow = []string{}
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)

	renderedDailyView := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)

	comboView := lipgloss.JoinVertical(lipgloss.Left, renderedDailyView, m.helpView())
	return comboView
}
