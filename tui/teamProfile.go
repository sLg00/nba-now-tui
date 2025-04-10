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
	width            int
	height           int
	mainPort         viewport.Model
	tables           []table.Model
	tableNames       []string
	activeTableIndex int
	quitting         bool
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
		mainPort:         vp,
		width:            size.Width,
		height:           size.Height,
		tables:           make([]table.Model, 3),
		tableNames:       []string{"Team Info", "SEASON STATS", "ROSTER"},
		activeTableIndex: 1,
		quitting:         false,
	}

	cmds := tea.Batch(fetchBasicTeamInfoMsg(teamID),
		fetchTeamSeasonSnapshotMsg(teamID),
		fetchPlayerIndexMsg(teamID))

	return m, cmds, nil
}

// fetchBasicTeamInfoMsg is the exception to the rule that all TUI tables should use the generic buildTables function.
// The reason is that this function builds a unique, two column, single row table which holds the team logo.
// So it does not conform to the table standards.
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

func (m *TeamProfile) assembleTables() {
	if len(m.tables) == 0 {
		return
	}

	//the team info table is never focused
	m.tables[0] = m.tables[0].Focused(false)

	for i := 1; i < len(m.tables); i++ {
		if i == m.activeTableIndex {
			m.tables[i] = m.tables[i].Focused(true)
		} else {
			m.tables[i] = m.tables[i].Focused(false)
		}
	}

	centered := CenterStyle(m.mainPort.Width - 4)
	headerStyle := lipgloss.NewStyle().Bold(true)
	activeHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5")) //TODO: team color

	var sections []string
	sections = append(sections, InvisibleTableStyle.Render(m.tables[0].View()), "\n\n")

	for i := 1; i < len(m.tables); i++ {
		var headerContent string
		if i == m.activeTableIndex {
			headerContent = activeHeaderStyle.Render(" << " + m.tableNames[i] + " >> ")
		} else {
			headerContent = headerStyle.Render(" << " + m.tableNames[i] + " >> ")
		}

		sections = append(sections,
			centered.Render(headerContent),
			centered.Render(TableStyle.Render(m.tables[i].View())),
			"\n\n")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	m.mainPort.SetContent(content)
}

func (m *TeamProfile) Init() tea.Cmd { return nil }

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
		m.tables[0] = msg.teamBasicInfo
		m.assembleTables()
		return m, nil
	case teamSeasonSnapshotFetchedMsg:
		if msg.err != nil {
			log.Println("could not load season snapshot:", msg.err)
			return m, nil
		}
		m.tables[1] = msg.teamSeasonSnapshot
		m.assembleTables()
		return m, nil
	case playerIndexFetchedMsg:
		if msg.err != nil {
			log.Println("could not load player index:", msg.err)
			return m, nil
		}
		m.tables[2] = msg.roster
		m.assembleTables()
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Tab):
			if len(m.tables) > 1 {
				m.activeTableIndex = (m.activeTableIndex + 1) % len(m.tables)
				if m.activeTableIndex == 0 {
					m.activeTableIndex = 1
				}
				m.assembleTables()
			}
		case key.Matches(msg, Keymap.Back):
			ss, cmd, _ := NewSeasonStandings(WindowSize)
			return ss, cmd
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.mainPort.Width = msg.Width - 4
		m.mainPort.Height = msg.Height - 8
		m.width = msg.Width
		m.height = msg.Height
		m.assembleTables()
	}

	m.mainPort, cmd = m.mainPort.Update(msg)
	cmds = append(cmds, cmd)

	if len(m.tables) > m.activeTableIndex {
		var tableCmd tea.Cmd
		m.tables[m.activeTableIndex], tableCmd = m.tables[m.activeTableIndex].Update(msg)
		cmds = append(cmds, tableCmd)
		m.assembleTables()
	}

	return m, tea.Batch(cmds...)
}

func (m *TeamProfile) helpView() string {

	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m *TeamProfile) View() string {
	if m.quitting {
		return ""
	}

	if len(m.tables) == 0 {
		return DocStyle.Render("Loading team profile data...")
	}

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
