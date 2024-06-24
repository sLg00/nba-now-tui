package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sLg00/nba-now-tui/app/datamodels"
	"log"
)

type seasonStandings struct {
	quitting  bool
	westTeams table.Model
	eastTeams table.Model
}

func (m seasonStandings) Init() tea.Cmd { return nil }

func initSeasonStandings(i list.Item, p *tea.Program) (*seasonStandings, error) {
	teams, headers, err := datamodels.PopulateTeamStats()
	if err != nil {
		log.Println("Could not populate player stats, error:", err)
		return nil, err
	}

	eastTeams, westTeams := teams.SplitStandingsPerConference()

	eastTeamsStrings := datamodels.ConvertToString(eastTeams)
	westTeamsStrings := datamodels.ConvertToString(westTeams)

	var (
		columns  []table.Column
		column   table.Column
		rows     []table.Row
		row      table.Row
		westRows []table.Row
		westRow  table.Row
	)

	for _, h := range headers {
		column = table.Column{
			Title: h,
			Width: 10,
		}
		columns = append(columns, column)
	}

	for _, r := range eastTeamsStrings {
		row = r
		rows = append(rows, row)
	}

	for _, r := range westTeamsStrings {
		westRow = r
		westRows = append(westRows, westRow)
	}

	eastTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15))

	westTable := table.New(
		table.WithColumns(columns),
		table.WithRows(westRows),
		table.WithFocused(false),
		table.WithHeight(15))

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	eastTable.SetStyles(s)
	westTable.SetStyles(s)
	m := &seasonStandings{eastTeams: eastTable, westTeams: westTable}

	return m, nil
}

func (m seasonStandings) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
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
	}
	m.eastTeams, cmd = m.eastTeams.Update(msg)
	m.westTeams, cmd = m.westTeams.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m seasonStandings) helpView() string {
	// TODO: use the keymaps to populate the help string
	return HelpStyle("\n ↑/↓: navigate  • backspace: back • q: quit\n")
}

func (m seasonStandings) View() string {
	if m.quitting {
		return ""
	}
	renderedEastTable := TableStyle.Render(m.eastTeams.View()) + "\n"
	renderedWestTable := TableStyle.Render(m.westTeams.View()) + "\n"
	comboView := lipgloss.JoinVertical(lipgloss.Left, "\n", renderedEastTable, renderedWestTable, m.helpView())
	return DocStyle.Render(comboView)
}
