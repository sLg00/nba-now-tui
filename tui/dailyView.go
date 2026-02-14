package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/cmd/converters"
	filesystemops "github.com/sLg00/nba-now-tui/cmd/nba/filesystem"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"github.com/sLg00/nba-now-tui/cmd/nba/pathManager"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

type focusZone int

const (
	focusDateSelector focusZone = iota
	focusGameCards
)

type DailyView struct {
	dateSelector DateSelector
	gameCards    []table.Model
	quitting     bool
	focusIndex   int
	numCols      int
	width        int
	height       int
	focus        focusZone
	loading      bool
}

type dailyScoresFetchedMsg struct {
	scores [][]string
	err    error
}

type dateScoresFetchedMsg struct {
	scores [][]string
	date   string
	err    error
}

type gameDataFetchedMsg struct {
	err error
}

func NewDailyView(size tea.WindowSizeMsg) (*DailyView, tea.Cmd, error) {
	client := nbaAPI.NewClient()
	currentDate, err := client.Dates.GetCurrentDate()
	if err != nil {
		return &DailyView{}, nil, fmt.Errorf("failed to get current date: %w", err)
	}

	ds := NewDateSelector(currentDate)
	ds.SetWidth(size.Width)

	m := &DailyView{
		dateSelector: ds,
		focusIndex:   0,
		numCols:      3,
		width:        size.Width,
		height:       size.Height,
		focus:        focusDateSelector,
	}

	cl, err := client.Loader.LoadDailyScoreboard()
	_, _, err = converters.PopulateDailyGameResults(cl)
	if err != nil {
		return &DailyView{}, nil, fmt.Errorf("failed to populate daily scores: %w\n", err)
	}

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
		cl, err := nbaAPI.NewClient().Loader.LoadDailyScoreboard()
		if err != nil {
			log.Println("error loading daily scoreboard:", err)
		}
		dailyScores, _, err := converters.PopulateDailyGameResults(cl)
		if err != nil {
			return dailyScoresFetchedMsg{err: err}
		}
		return dailyScoresFetchedMsg{scores: types.ConvertToStringMatrix(dailyScores)}
	}
}

// fetchGameDataCmd first checks the game's status (1,2,3), then queries the NBA API for box score results if
// the status is greater than 1
func fetchGameDataCmd(gameID string) tea.Cmd {
	return func() tea.Msg {
		gameStatus, err := converters.CheckGameStatus(gameID)
		if err != nil {
			return gameDataFetchedMsg{err: err}
		}
		if gameStatus > 1 {
			err = nbaAPI.NewClient().FetchBoxScore(gameID)
			return gameDataFetchedMsg{err: err}
		}

		return gameDataFetchedMsg{err: err}
	}
}

func fetchDailyScoresForDateCmd(date string) tea.Cmd {
	return func() tea.Msg {
		err := nbaAPI.NewClient().FetchDailyScoresForDate(date)
		if err != nil {
			return dateScoresFetchedMsg{date: date, err: err}
		}

		paths := pathManager.PathFactoryForDate(date)
		loader := filesystemops.NewDataLoader(filesystemops.NewDefaultFsHandler(), paths)
		cl, err := loader.LoadDailyScoreboard()
		if err != nil {
			return dateScoresFetchedMsg{date: date, err: err}
		}

		dailyScores, _, err := converters.PopulateDailyGameResults(cl)
		if err != nil {
			return dateScoresFetchedMsg{date: date, err: err}
		}

		return dateScoresFetchedMsg{
			scores: types.ConvertToStringMatrix(dailyScores),
			date:   date,
		}
	}
}

// getGameId extrapolates the gameID from gameCard in order to query the NBA API for the corresponding box score
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
		m.gameCards = nil
		for _, score := range msg.scores {
			gameStatus, err := converters.CheckGameStatus(score[0])
			if err != nil {
				log.Println("Error while querying game status", err)
			}
			if gameStatus > 0 {
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
		}
		m.focusIndex = 0

	case dateScoresFetchedMsg:
		m.loading = false
		if msg.err != nil {
			log.Println("Error fetching scores for date", msg.date, msg.err)
			return m, nil
		}
		m.gameCards = nil
		for _, score := range msg.scores {
			gameStatus, err := converters.CheckGameStatus(score[0])
			if err != nil {
				log.Println("Error while querying game status", err)
			}
			if gameStatus > 0 {
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
		}
		m.focusIndex = 0

	case dateChangedMsg:
		m.loading = true
		m.gameCards = nil
		return m, fetchDailyScoresForDateCmd(msg.date)

	case gameDataFetchedMsg:
		if msg.err != nil {
			log.Println("Error while fetching game data", msg.err)
		}

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
		}

		if m.focus == focusDateSelector {
			switch {
			case key.Matches(msg, Keymap.Tab):
				if len(m.gameCards) > 0 {
					m.dateSelector.Blur()
					m.focus = focusGameCards
				}
				return m, nil
			case key.Matches(msg, Keymap.Down):
				if !m.dateSelector.editing && len(m.gameCards) > 0 {
					m.dateSelector.Blur()
					m.focus = focusGameCards
					return m, nil
				}
			}

			var dsCmd tea.Cmd
			m.dateSelector, dsCmd = m.dateSelector.Update(msg)
			cmds = append(cmds, dsCmd)
			return m, tea.Batch(cmds...)
		}

		switch {
		case key.Matches(msg, Keymap.Tab):
			m.focus = focusDateSelector
			m.dateSelector.Focus()
			return m, nil
		case key.Matches(msg, Keymap.Enter):
			gameID, err := m.getGameId()
			if err != nil {
				log.Printf("Could not get game id: %v", err)
			}
			gameStatus, err := converters.CheckGameStatus(gameID)
			if err != nil {
				log.Println("Error while querying game status", err)
			}
			if gameStatus > 1 {
				bx, cmd, err := NewBoxScore(gameID, WindowSize)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
				return bx, cmd
			}
		case key.Matches(msg, Keymap.Up):
			if m.focusIndex < m.numCols {
				m.focus = focusDateSelector
				m.dateSelector.Focus()
				return m, nil
			}
			m.focusIndex -= m.numCols
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

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.dateSelector.SetWidth(msg.Width)
	}

	if len(m.gameCards) > 0 && m.focus == focusGameCards {
		m.gameCards[m.focusIndex], cmd = m.gameCards[m.focusIndex].Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func renderDailyView(m DailyView) string {
	var content string
	var rows []string
	var currentRow []string

	dateView := m.dateSelector.View()

	if m.loading {
		content = dateView + "\n\nLoading..."
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render(content)
	}

	for i, gameCard := range m.gameCards {
		if i == m.focusIndex && m.focus == focusGameCards {
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

	if len(m.gameCards) == 0 {
		content = dateView + "\n\n" + "No games happened during " + m.dateSelector.date
	} else {
		if len(currentRow) > 0 {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
		}
		cardsContent := lipgloss.JoinVertical(lipgloss.Left, rows...)
		content = dateView + "\n" + cardsContent
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
