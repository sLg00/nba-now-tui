package tui

import (
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
	teamTable       table.Model
	width           int
	height          int
	mainPort        viewport.Model
	teamBasicInfo   table.Model
	seasonStatsPort table.Model
	rosterPort      table.Model
}

type teamBasicInfoFetchedMsg struct {
	err           error
	teamBasicInfo table.Model
}

func NewTeamProfile(teamID string, size tea.WindowSizeMsg) (*TeamProfile, tea.Cmd, error) {
	m := &TeamProfile{
		width:  size.Width,
		height: size.Height,
	}

	cl, err := nbaAPI.NewClient().Loader.LoadTeamInfo(teamID)
	_, _, err = converters.PopulateTeamInfo(cl)
	if err != nil {
		return &TeamProfile{}, nil, err
	}

	cmd := fetchBasicTeamInfoMsg(teamID)

	return m, cmd, nil
}

func fetchBasicTeamInfoMsg(teamID string) tea.Cmd {
	return func() tea.Msg {
		cl, err := nbaAPI.NewClient().Loader.LoadTeamInfo(teamID)
		if err != nil {
			log.Println("error loading team profile:", err)
		}
		data, _, err := converters.PopulateTeamInfo(cl)
		if err != nil {
			return teamBasicInfoFetchedMsg{err: err}
		}

		teamBasicsStrings := types.ConvertToStringFlat(data)
		season := teamBasicsStrings[1]
		name := teamBasicsStrings[3]
		//TODO: remove this stupid temporary hack
		if name == "Trail Blazers" {
			name = "TrailBlazers"
		}

		city := teamBasicsStrings[2]
		conf := teamBasicsStrings[5]
		div := teamBasicsStrings[6]
		//wins := teamBasicsStrings[9]
		//losses := teamBasicsStrings[10]
		//winPct := teamBasicsStrings[11]
		//confRank := teamBasicsStrings[12]
		//divRank := teamBasicsStrings[13]

		logo := logos.LoadTeamLogo(name)

		var columns []table.Column
		var rows []table.Row
		logoColumn := table.NewColumn("logo", "", 130)
		dataColumn := table.NewColumn("data", "", 100)
		columns = append(columns, logoColumn, dataColumn)

		rowData := make(table.RowData)
		rowData["logo"] = logo
		rowData["data"] = city + " " + name + "\n\n" +
			season + " Season\n\n " + conf + " | " + div + "\n\n"

		row := table.NewRow(rowData)
		rows = append(rows, row)

		basicInfoTable := table.New(columns).WithRows(rows).
			WithHeaderVisibility(false).
			WithMultiline(true).WithBaseStyle(InvisibleTableStyle)
		return teamBasicInfoFetchedMsg{err: nil, teamBasicInfo: basicInfoTable}
	}
}

func (m TeamProfile) Init() tea.Cmd { return nil }

func (m TeamProfile) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
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
		m := &TeamProfile{teamBasicInfo: msg.teamBasicInfo}
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Back):
			return InitMenu()
		case key.Matches(msg, Keymap.Quit):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.teamBasicInfo, cmd = m.teamBasicInfo.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m TeamProfile) helpView() string {

	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m TeamProfile) View() string {
	renderedTable := InvisibleTableStyle.Render(m.teamBasicInfo.View()) + "\n"
	comboView := lipgloss.JoinVertical(lipgloss.Left, "\n", renderedTable, m.helpView())
	return DocStyle.Render(comboView)
}
