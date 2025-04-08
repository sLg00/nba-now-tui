package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/assets/logos"
	"github.com/sLg00/nba-now-tui/cmd/converters"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"log"
)

type TeamProfile struct {
	width              int
	height             int
	mainPort           viewport.Model
	teamBasicInfo      table.Model
	teamSeasonSnapshot table.Model
	seasonStatsPort    table.Model
	rosterPort         table.Model
	activeTable        int
	focused            bool
}

type teamBasicInfoFetchedMsg struct {
	err           error
	teamBasicInfo table.Model
}

type teamSeasonSnapshotFetchedMsg struct {
	err                error
	teamSeasonSnapshot table.Model
}

type playerIndexFetchedMsg struct {
	err    error
	roster table.Model
}

func NewTeamProfile(teamID string, size tea.WindowSizeMsg) (*TeamProfile, tea.Cmd, error) {
	vp := viewport.New(size.Width-4, size.Height-8)

	dataStr, _, err := TeamDataStrings(teamID)
	if err != nil {
		return nil, nil, err
	}
	name := dataStr[3]

	teamStyle := TeamViewPortStyle(TeamColor(name))
	vp.Style = teamStyle

	m := &TeamProfile{
		mainPort: vp,
		width:    size.Width,
		height:   size.Height,
	}

	cmds := tea.Batch(fetchBasicTeamInfoMsg(teamID), fetchTeamSeasonSnapshotMsg(teamID), fetchPlayerIndexMsg(teamID))

	return m, cmds, nil
}

func fetchBasicTeamInfoMsg(teamID string) tea.Cmd {
	return func() tea.Msg {
		teamBasicsStrings, _, err := TeamDataStrings(teamID)
		if err != nil {
			return teamBasicInfoFetchedMsg{err: err}
		}

		season := teamBasicsStrings[1]
		name := teamBasicsStrings[3]
		//TODO: remove this stupid temporary hack
		if name == "Trail Blazers" {
			name = "TrailBlazers"
		}

		city := teamBasicsStrings[2]
		conf := teamBasicsStrings[5]
		div := teamBasicsStrings[6]

		logo := logos.LoadTeamLogo(name)

		var columns []table.Column
		var rows []table.Row
		logoColumn := table.NewColumn("logo", "", 130)
		dataColumn := table.NewColumn("data", "", WindowSize.Width-142)
		columns = append(columns, logoColumn, dataColumn)

		rowData := make(table.RowData)
		rowData["logo"] = logo
		rowData["data"] = city + " " + name + "\n\n" +
			season + " Season\n\n " + conf + " | " + div + "\n\n"

		row := table.NewRow(rowData)
		rows = append(rows, row)

		basicInfoTable := table.New(columns).WithRows(rows).
			WithHeaderVisibility(false).
			WithMultiline(true).WithBaseStyle(InvisibleTableStyle).Focused(false)
		return teamBasicInfoFetchedMsg{err: nil, teamBasicInfo: basicInfoTable}
	}
}

func fetchTeamSeasonSnapshotMsg(teamID string) tea.Cmd {
	return func() tea.Msg {
		teamBasicsStrings, headers, err := TeamDataStrings(teamID)
		if err != nil {
			return teamBasicInfoFetchedMsg{err: err}
		}

		seasonSnapsShotTable := buildTables(headers, teamBasicsStrings, types.TeamCommonInfo{})

		return teamSeasonSnapshotFetchedMsg{err: nil, teamSeasonSnapshot: seasonSnapsShotTable}

	}
}

func fetchPlayerIndexMsg(teamID string) tea.Cmd {
	return func() tea.Msg {
		cl, err := nbaAPI.NewClient().Loader.LoadPlayerIndex(teamID)
		if err != nil {
			return playerIndexFetchedMsg{err: err}
		}

		players, headers, err := converters.PopulatePlayerIndex(cl)
		if err != nil {
			return playerIndexFetchedMsg{err: err}
		}

		playerStrings := types.ConvertToStringMatrix(players)
		tableModel := buildTables(headers, playerStrings, types.IndexPlayer{}).Focused(true)

		return playerIndexFetchedMsg{roster: tableModel, err: nil}
	}
}

func (m *TeamProfile) Init() tea.Cmd { return nil }

func (m *TeamProfile) updateViewPortContent() {

	centered := CenterStyle(m.mainPort.Width - 4)

	sections := []string{
		InvisibleTableStyle.Render(m.teamBasicInfo.View()),
		"\n\n",
		centered.Render(" << SEASON AT A GLANCE >>"),
		centered.Render(m.teamSeasonSnapshot.View()),
		"\n\n",
		centered.Render(" << ROSTER >>"),
		centered.Render(m.rosterPort.View()),
		"\n\n",
		centered.Render("<< MOAR >>"),
	}
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	m.mainPort.SetContent(content)
}

func (m *TeamProfile) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case teamProfileDownloadedMsg:
		if msg.err != nil {
			log.Println("could not download team profile:", msg.err)
			return m, nil
		}
	case teamBasicInfoFetchedMsg:
		if msg.err != nil {
			log.Println("could not load team profile:", msg.err)
			return m, nil
		}
		m.teamBasicInfo = msg.teamBasicInfo
		m.updateViewPortContent()
		return m, nil
	case teamSeasonSnapshotFetchedMsg:
		if msg.err != nil {
			log.Println("could not load season snapshot:", msg.err)
			return m, nil
		}
		m.teamSeasonSnapshot = msg.teamSeasonSnapshot
		m.updateViewPortContent()
		return m, nil
	case playerIndexFetchedMsg:
		if msg.err != nil {
			log.Println("could not load player index:", msg.err)
			return m, nil
		}
		m.rosterPort = msg.roster
		m.updateViewPortContent()
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Tab):

		case key.Matches(msg, Keymap.Back):
			ss, cmd, _ := NewSeasonStandings(WindowSize)
			return ss, cmd
		case key.Matches(msg, Keymap.Quit):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.mainPort.Width = msg.Width - 2
		m.mainPort.Height = msg.Height - 2
	}

	m.mainPort, cmd = m.mainPort.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *TeamProfile) helpView() string {

	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m *TeamProfile) View() string {
	comboView := lipgloss.JoinVertical(lipgloss.Left, m.mainPort.View(), m.helpView())
	return DocStyle.Render(comboView)
}

// TeamDataStrings is a helper function to reduce code duplication (just converts the structs to strings)
func TeamDataStrings(teamID string) ([]string, []string, error) {
	cl, err := nbaAPI.NewClient().Loader.LoadTeamInfo(teamID)
	data, headers, err := converters.PopulateTeamInfo(cl)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting team info: %w", err)
	}

	dataStr := types.ConvertToStringFlat(data)

	return dataStr, headers, nil
}
