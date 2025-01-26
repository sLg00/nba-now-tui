package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/cmd/client"
	"github.com/sLg00/nba-now-tui/cmd/datamodels"
	"log"
	"os"
)

type DailyView struct {
	gameCards  []table.Model
	quitting   bool
	focusIndex int
	numCols    int
	width      int
	height     int
}

type dailyScoresFetchedMsg struct {
	scores [][]string
	err    error
}

type gameDataFetchedMsg struct {
	err error
}

func NewDailyView(size tea.WindowSizeMsg) (*DailyView, tea.Cmd, error) {
	m := &DailyView{
		focusIndex: 0,
		numCols:    3,
		width:      size.Width,
		height:     size.Height,
	}

	// Attempt to fetch initial data directly to validate API availability
	_, _, err := datamodels.PopulateDailyGameResults(datamodels.UnmarshallResponseJSON)
	if err != nil {
		return &DailyView{}, nil, fmt.Errorf("failed to populate daily scores: %w", err)
	}

	// Prepare the Init command for fetching dynamic updates
	cmd := fetchDailyScoresCmd()

	return m, cmd, nil
}

func newGameCard(r []table.Row) (table.Model, error) {
	columns := []table.Column{
		table.NewColumn("teams", "", 7),
		table.NewColumn("scores", "", 7),
	}

	gc := table.New(columns)
	if len(r) > 0 {
		gc = gc.WithRows(r).WithBaseStyle(lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder()).
			BorderForeground(lipgloss.Color("7"))).
			WithHeaderVisibility(false)
		gc = gc.BorderRounded()
		return gc, nil
	}
	err := fmt.Errorf("no games happened during the date")
	return table.Model{}, err
}

func (m DailyView) Init() tea.Cmd { return fetchDailyScoresCmd() }

func fetchDailyScoresCmd() tea.Cmd {
	return func() tea.Msg {
		dailyScores, _, err := datamodels.PopulateDailyGameResults(datamodels.UnmarshallResponseJSON)
		if err != nil {
			return dailyScoresFetchedMsg{err: err}
		}
		return dailyScoresFetchedMsg{scores: datamodels.ConvertToStringMatrix(dailyScores)}
	}
}

func fetchGameDataCmd(gameID string) tea.Cmd {
	return func() tea.Msg {
		err := client.NewClient().MakeOnDemandRequests(gameID)
		return gameDataFetchedMsg{err: err}
	}
}

func (m DailyView) getGameId() (string, error) {
	focusedCard := m.gameCards[m.focusIndex]
	rows := focusedCard.GetVisibleRows()
	if len(rows) > 0 {
		gameId, ok := rows[0].Data["gameID"].(string)
		if ok {
			return gameId, nil
		}
	}
	if len(rows) > 1 || len(rows) < 1 {
		log.Println("Either 0 rows or more than 1 row were selected")
		err := fmt.Errorf("game ID not found")
		return "", err
		//TODO: Display pop-up with User error! :)
	}
	return "", nil
}

func (m DailyView) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case dailyScoresFetchedMsg:
		if msg.err != nil {
			log.Println("Error while fetching daily game data", msg.err)
			return m, nil
		}
		for _, score := range msg.scores {
			rows := []table.Row{
				table.NewRow(table.RowData{"teams": score[4], "scores": score[3], "gameID": score[0]}),
				table.NewRow(table.RowData{"teams": score[8], "scores": score[7], "gameID": score[0]}),
			}
			gameCard, err := newGameCard(rows)
			if err != nil {
				log.Println("Error creating game card", err)
				continue
			}
			m.gameCards = append(m.gameCards, gameCard)
			cmds = append(cmds, fetchGameDataCmd(score[0]))
		}
	case gameDataFetchedMsg:
		if msg.err != nil {
			log.Println("Error while fetching game data", msg.err)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Back):
			return InitMenu()
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, Keymap.Enter):
			gameID, err := m.getGameId()
			if err != nil {
				log.Printf("Could not get game id: %v", err)
			}
			bx, cmd, err := NewBoxScore(gameID, WindowSize)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			return bx, cmd
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

	if len(m.gameCards) > 0 {
		m.gameCards[m.focusIndex], cmd = m.gameCards[m.focusIndex].Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}
	return m, tea.Batch(cmds...)
}

func renderDailyView(m DailyView) string {
	var content string
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

	dateToDisplayInCaseOfEmpty, _ := client.GetDateArg()
	if len(m.gameCards) == 0 {
		content = "No games happened during " + dateToDisplayInCaseOfEmpty
	} else {
		if len(currentRow) > 0 {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))

		}
		content = lipgloss.JoinVertical(lipgloss.Left, rows...)
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}

func (m DailyView) renderHelpView() string {
	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m DailyView) View() string {
	if m.quitting {
		return ""
	}
	comboView := lipgloss.JoinVertical(lipgloss.Left, renderDailyView(m), m.renderHelpView())
	return comboView
}
