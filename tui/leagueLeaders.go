package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/app/datamodels"
	"log"
)

// removeIndex is a helper function to delete elements from a slice
func removeIndex[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

type leagueLeaders struct {
	leaderboard table.Model
	quitting    bool
	height      int
	width       int
	maxHeight   int
	maxWidth    int
}

func (m leagueLeaders) Init() tea.Cmd { return nil }

// initLeagueLeadersTable returns a table.Model which is populated with the current league leaders (per PPG)
func initLeagueLeaders(i list.Item, p *tea.Program) (*leagueLeaders, error) {
	playerStats, headers, err := datamodels.PopulatePlayerStats()
	log.Println("MINGE MUNNI")
	if err != nil {
		log.Println("Could not populate player stats, error:", err)
		return nil, err
	}
	playerStatsString := datamodels.ConvertToString(playerStats)

	var filteredHeaders []string
	for _, h := range headers {
		if !strings.Contains(h, "ID") {
			filteredHeaders = append(filteredHeaders, h)
		}
	}

	var filteredRows [][]string
	for _, fr := range playerStatsString {
		fr = slices.Delete(fr, 0, 1)
		fr = slices.Delete(fr, 2, 3)
		filteredRows = append(filteredRows, fr)
	}

	var sortedHeaders []string
	sortedHeaders = append(sortedHeaders, filteredHeaders[0:5]...)
	sortedHeaders = append(sortedHeaders, filteredHeaders[21:23]...)
	sortedHeaders = append(sortedHeaders, filteredHeaders[17:21]...)
	sortedHeaders = append(sortedHeaders, filteredHeaders[5:17]...)

	var sortedRows [][]string
	for _, r := range filteredRows {
		var sorterRow []string
		sorterRow = append(sorterRow, r[0:5]...)
		sorterRow = append(sorterRow, r[21:23]...)
		sorterRow = append(sorterRow, r[17:21]...)
		sorterRow = append(sorterRow, r[5:17]...)
		sortedRows = append(sortedRows, sorterRow)
	}

	var (
		columns []table.Column
		column  table.Column
		rows    []table.Row
		row     table.Row
	)

	for _, h := range sortedHeaders {
		column = table.NewColumn(h, h, 15)
		columns = append(columns, column)

	}

	for _, r := range sortedRows {
		rowData := make(table.RowData)
		for i, rd := range r {
			columnTitle := columns[i].Title()
			rowData[columnTitle] = rd
		}
		row = table.NewRow(rowData)
		rows = append(rows, row)
	}

	t := table.New(columns).
		WithRows(rows).
		WithHeaderVisibility(true).
		Focused(true).
		WithPageSize(20).
		WithMaxTotalWidth(150).
		WithHorizontalFreezeColumnCount(2).WithBaseStyle(lipgloss.NewStyle().BorderStyle(lipgloss.DoubleBorder()))

	m := &leagueLeaders{leaderboard: t, maxHeight: 25, maxWidth: 80}

	return m, nil
}

func (m leagueLeaders) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Back):
			return InitMenu()
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
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
	}

	m.leaderboard, cmd = m.leaderboard.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m leagueLeaders) helpView() string {
	// TODO: use the keymaps to populate the help string
	return HelpStyle("\n ↑/↓: navigate  • backspace: back • q: quit\n")
}

func (m leagueLeaders) View() string {
	if m.quitting {
		return ""
	}
	renderedLeaders := TableStyle.Render(m.leaderboard.View()) + "\n"
	comboView := lipgloss.JoinVertical(lipgloss.Left, "\n", renderedLeaders, m.helpView())
	return DocStyle.Render(comboView)
}
