package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/cmd/converters"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"log"
)

type LeagueLeaders struct {
	leaderboard table.Model
	quitting    bool
	height      int
	width       int
	maxHeight   int
	maxWidth    int
}

type fetchLeagueLeadersMsg struct {
	err   error
	table table.Model
}

func NewLeagueLeaders(size tea.WindowSizeMsg) (*LeagueLeaders, tea.Cmd, error) {
	m := &LeagueLeaders{
		height: size.Height,
		width:  size.Width,
	}

	cl, err := nbaAPI.NewClient().Loader.LoadLeagueLeaders()
	_, _, err = converters.PopulatePlayerStats(cl)
	if err != nil {
		return &LeagueLeaders{}, nil, fmt.Errorf("failed to populate player stats: %w", err)
	}

	cmd := fetchLeagueLeadersCmd()

	return m, cmd, nil
}

// fetchLeagueLeadersCmd fetches the required data from pre-existing JSON files and creates the table structure and rows
func fetchLeagueLeadersCmd() tea.Cmd {
	return func() tea.Msg {
		cl, err := nbaAPI.NewClient().Loader.LoadLeagueLeaders()
		if err != nil {
			log.Println("failed to load league leaders")
		}
		playerStats, headers, err := converters.PopulatePlayerStats(cl)
		if err != nil {
			return fetchLeagueLeadersMsg{err: err}
		}

		playerStatsString := types.ConvertToStringMatrix(playerStats)
		tableModel := buildTables(headers, playerStatsString, types.Player{}).
			Focused(true).
			WithBaseStyle(TableStyle).
			WithPageSize(20)

		return fetchLeagueLeadersMsg{table: tableModel, err: nil}
	}
}

func (m LeagueLeaders) Init() tea.Cmd { return nil }

func (m LeagueLeaders) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case fetchLeagueLeadersMsg:
		if msg.err != nil {
			return m, nil
		}
		m := &LeagueLeaders{leaderboard: msg.table, maxHeight: 25, maxWidth: 125}
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Back):
			return InitMenu()
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, Keymap.Enter):
			selectedRows := m.leaderboard.SelectedRows()
			if len(selectedRows) == 1 {
				playerID := selectedRows[0].Data["PLAYER_ID"].(string)
				log.Println(playerID)
				//TODO: add player profile init logic
			}
			if len(selectedRows) > 1 || len(selectedRows) < 1 {
				log.Println("Either 0 rows or more than 1 row were selected")
				//TODO: Display pop-up with User error! :)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if m.width > m.maxWidth {
			m.maxWidth = m.width
		}
		if m.height > m.maxHeight {
			m.maxHeight = m.height
		}

		pageSize := calculatePageSize(msg.Height, 1)
		m.leaderboard = m.leaderboard.WithPageSize(pageSize).WithFooterVisibility(false)
	}

	m.leaderboard, cmd = m.leaderboard.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m LeagueLeaders) helpView() string {
	return HelpStyle(HelpFooter())
}

func (m LeagueLeaders) View() string {
	if m.quitting {
		return ""
	}
	comboView := lipgloss.JoinVertical(lipgloss.Left,
		m.leaderboard.View(),
		m.helpView())
	return DocStyle.Render(comboView)
}
