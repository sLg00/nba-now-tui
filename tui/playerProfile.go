package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sLg00/nba-now-tui/cmd/converters"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"log"
	"strconv"
)

type PlayerProfile struct {
	width            int
	height           int
	mainPort         viewport.Model
	bio              *types.PlayerBio
	currentStats     *types.SeasonStats
	tables           []table.Model
	tableNames       []string
	activeTableIndex int
	backView         string
	sourceDate       string
	teamColor        lipgloss.Color
	quitting         bool
}

type playerProfileDownloadedMsg struct {
	err        error
	playerID   string
	backView   string
	sourceDate string
}

func downloadPlayerProfile(playerID string, backView string, sourceDate string) tea.Cmd {
	return func() tea.Msg {
		err := nbaAPI.NewClient().FetchPlayerProfile(playerID)
		return playerProfileDownloadedMsg{err: err, playerID: playerID, backView: backView, sourceDate: sourceDate}
	}
}

type playerBioFetchedMsg struct {
	err error
	bio types.PlayerBio
}

type playerCareerStatsFetchedMsg struct {
	err          error
	seasonStats  table.Model
	currentStats *types.SeasonStats
}

type playerGameLogFetchedMsg struct {
	err     error
	gameLog table.Model
}

func NewPlayerProfile(playerID string, backView string, sourceDate string, size tea.WindowSizeMsg) (*PlayerProfile, tea.Cmd, error) {
	vp := viewport.New(size.Width-4, size.Height-8)
	vp.Style = TeamViewPortStyle(lipgloss.Color("#FFFFFF"))

	m := &PlayerProfile{
		mainPort:         vp,
		width:            size.Width,
		height:           size.Height,
		tables:           make([]table.Model, 2),
		tableNames:       []string{"LAST 5 GAMES", "CAREER STATS"},
		activeTableIndex: 0,
		backView:         backView,
		sourceDate:       sourceDate,
		teamColor:        lipgloss.Color("#FFFFFF"),
		quitting:         false,
	}

	cmds := tea.Batch(
		fetchPlayerBioCmd(playerID),
		fetchPlayerCareerStatsCmd(playerID),
		fetchPlayerGameLogCmd(playerID),
	)

	return m, cmds, nil
}

func fetchPlayerBioCmd(playerID string) tea.Cmd {
	return func() tea.Msg {
		cl := nbaAPI.NewClient()
		rs, err := cl.Loader.LoadPlayerInfo(playerID)
		if err != nil {
			return playerBioFetchedMsg{err: err}
		}

		bio, err := converters.PopulatePlayerBio(rs)
		if err != nil {
			return playerBioFetchedMsg{err: err}
		}

		return playerBioFetchedMsg{bio: bio}
	}
}

func fetchPlayerCareerStatsCmd(playerID string) tea.Cmd {
	return func() tea.Msg {
		cl := nbaAPI.NewClient()
		rs, err := cl.Loader.LoadPlayerCareerStats(playerID)
		if err != nil {
			return playerCareerStatsFetchedMsg{err: err}
		}

		stats, headers, err := converters.PopulateSeasonStats(rs)
		if err != nil {
			return playerCareerStatsFetchedMsg{err: err}
		}

		var currentStats *types.SeasonStats
		if len(stats) > 0 {
			currentStats = &stats[0]
		}

		stringMatrix := types.ConvertToStringMatrix(stats)
		tableModel := buildTables(headers, stringMatrix, types.SeasonStats{})

		return playerCareerStatsFetchedMsg{
			seasonStats:  tableModel,
			currentStats: currentStats,
		}
	}
}

func fetchPlayerGameLogCmd(playerID string) tea.Cmd {
	return func() tea.Msg {
		cl := nbaAPI.NewClient()
		rs, err := cl.Loader.LoadPlayerGameLog(playerID)
		if err != nil {
			return playerGameLogFetchedMsg{err: err}
		}

		entries, headers, err := converters.PopulateGameLog(rs)
		if err != nil {
			return playerGameLogFetchedMsg{err: err}
		}

		stringMatrix := types.ConvertToStringMatrix(entries)
		tableModel := buildTables(headers, stringMatrix, types.GameLogEntry{})

		return playerGameLogFetchedMsg{gameLog: tableModel}
	}
}

func renderStatCard(label, value string, clr lipgloss.Color) string {
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Align(lipgloss.Center).Width(10)
	valueStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Align(lipgloss.Center).Width(10)
	cardStyle := lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(clr).Padding(0, 1)
	content := lipgloss.JoinVertical(lipgloss.Center, labelStyle.Render(label), valueStyle.Render(value))
	return cardStyle.Render(content)
}

func (m *PlayerProfile) renderTopSection() string {
	if m.bio == nil {
		return ""
	}

	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Background(m.teamColor).Padding(0, 2)
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	name := nameStyle.Render(m.bio.DisplayName)
	info := infoStyle.Render(fmt.Sprintf("%s %s | #%s | %s",
		m.bio.TeamCity, m.bio.TeamName, m.bio.JerseyNumber, m.bio.Position))

	bioLine := infoStyle.Render(fmt.Sprintf("Height: %s  Weight: %slb  Country: %s  Draft: %s R%s Pick %s  School: %s  Experience: %d Years",
		m.bio.Height, m.bio.Weight, m.bio.Country,
		m.bio.DraftYear, m.bio.DraftRound, m.bio.DraftNumber,
		m.bio.School, m.bio.SeasonExp))

	sections := []string{name, info, ""}

	if m.currentStats != nil {
		cards := lipgloss.JoinHorizontal(lipgloss.Center,
			renderStatCard("GP", strconv.Itoa(m.currentStats.GP), m.teamColor),
			renderStatCard("PPG", strconv.FormatFloat(m.currentStats.PTS, 'f', 1, 64), m.teamColor),
			renderStatCard("RPG", strconv.FormatFloat(m.currentStats.REB, 'f', 1, 64), m.teamColor),
			renderStatCard("APG", strconv.FormatFloat(m.currentStats.AST, 'f', 1, 64), m.teamColor),
			renderStatCard("SPG", strconv.FormatFloat(m.currentStats.STL, 'f', 1, 64), m.teamColor),
			renderStatCard("BPG", strconv.FormatFloat(m.currentStats.BLK, 'f', 1, 64), m.teamColor),
		)
		sections = append(sections, cards, "")
	}

	sections = append(sections, bioLine)
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *PlayerProfile) assembleSections() {
	centered := CenterStyle(m.mainPort.Width - 4)
	headerStyle := lipgloss.NewStyle().Bold(true)
	activeHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))

	var sections []string

	topSection := m.renderTopSection()
	if topSection != "" {
		sections = append(sections, topSection, "\n")
	}

	for i := 0; i < len(m.tables); i++ {
		if i == m.activeTableIndex {
			m.tables[i] = m.tables[i].Focused(true)
		} else {
			m.tables[i] = m.tables[i].Focused(false)
		}

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

func (m *PlayerProfile) Init() tea.Cmd { return nil }

func (m *PlayerProfile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case playerBioFetchedMsg:
		if msg.err != nil {
			log.Println("could not load player bio:", msg.err)
			return m, nil
		}
		m.bio = &msg.bio
		teamName := msg.bio.TeamName
		m.teamColor = TeamColor(teamName)
		m.mainPort.Style = TeamViewPortStyle(m.teamColor)
		m.assembleSections()
		return m, nil

	case playerCareerStatsFetchedMsg:
		if msg.err != nil {
			log.Println("could not load career stats:", msg.err)
			return m, nil
		}
		m.currentStats = msg.currentStats
		m.tables[1] = msg.seasonStats
		m.assembleSections()
		return m, nil

	case playerGameLogFetchedMsg:
		if msg.err != nil {
			log.Println("could not load game log:", msg.err)
			return m, nil
		}
		m.tables[0] = msg.gameLog
		m.assembleSections()
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Tab):
			m.activeTableIndex = (m.activeTableIndex + 1) % len(m.tables)
			m.assembleSections()
		case key.Matches(msg, Keymap.Back):
			switch m.backView {
			case "boxscore":
				dv, cmd := NewDailyViewForDate(m.sourceDate, WindowSize)
				return dv, cmd
			case "leagueLeaders":
				ll, cmd, _ := NewLeagueLeaders(WindowSize)
				return ll, cmd
			default:
				return InitMenu()
			}
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.mainPort.Width = msg.Width - 4
		m.mainPort.Height = msg.Height - 8
		m.width = msg.Width
		m.height = msg.Height
		m.assembleSections()
	}

	var cmd tea.Cmd
	m.mainPort, cmd = m.mainPort.Update(msg)
	cmds = append(cmds, cmd)

	if len(m.tables) > m.activeTableIndex {
		var tableCmd tea.Cmd
		m.tables[m.activeTableIndex], tableCmd = m.tables[m.activeTableIndex].Update(msg)
		cmds = append(cmds, tableCmd)
		m.assembleSections()
	}

	return m, tea.Batch(cmds...)
}

func (m *PlayerProfile) helpView() string {
	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m *PlayerProfile) View() string {
	if m.quitting {
		return ""
	}

	if m.bio == nil {
		return DocStyle.Render("Loading player profile...")
	}

	comboView := lipgloss.JoinVertical(lipgloss.Left, m.mainPort.View(), m.helpView())
	return DocStyle.Render(comboView)
}
