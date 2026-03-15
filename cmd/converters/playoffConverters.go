package converters

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

// Series cursor index order (matches bracket renderer):
// 0-3:  East R1 (1v8, 4v5, 3v6, 2v7)
// 4-5:  East Semis
// 6:    East Finals
// 7:    NBA Finals
// 8:    West Finals
// 9-10: West Semis
// 11-14: West R1 (1v8, 4v5, 3v6, 2v7)

// ProjectedBracketFromStandings builds a pre-playoff bracket from current standings.
// All R2/CF/Finals series are TBD. First-round matchups follow 1v8, 4v5, 3v6, 2v7 seeding.
func ProjectedBracketFromStandings(east, west []types.Team, season string) types.PlayoffBracket {
	sort.Slice(east, func(i, j int) bool { return east[i].PlayoffSeeding < east[j].PlayoffSeeding })
	sort.Slice(west, func(i, j int) bool { return west[i].PlayoffSeeding < west[j].PlayoffSeeding })

	tbd := types.PlayoffTeam{Tricode: "TBD"}
	tbdSeries := func(round int, conf string) types.PlayoffSeries {
		return types.PlayoffSeries{
			Round: round, Conference: conf, Status: "pre",
			TopTeam: tbd, BottomTeam: tbd,
		}
	}

	teamFromStandings := func(t types.Team) types.PlayoffTeam {
		return types.PlayoffTeam{
			Name:    t.TeamName,
			Tricode: teamTricode(t),
			Seed:    t.PlayoffSeeding,
		}
	}

	r1Pairs := [][2]int{{0, 7}, {3, 4}, {2, 5}, {1, 6}} // seed indices: 1v8, 4v5, 3v6, 2v7

	eastR1 := make([]types.PlayoffSeries, 4)
	westR1 := make([]types.PlayoffSeries, 4)
	for i, pair := range r1Pairs {
		hi, lo := pair[0], pair[1]
		if lo < len(east) {
			eastR1[i] = types.PlayoffSeries{
				Round: 1, Conference: "East", Status: "pre",
				TopTeam: teamFromStandings(east[hi]), BottomTeam: teamFromStandings(east[lo]),
			}
		} else {
			eastR1[i] = tbdSeries(1, "East")
		}
		if lo < len(west) {
			westR1[i] = types.PlayoffSeries{
				Round: 1, Conference: "West", Status: "pre",
				TopTeam: teamFromStandings(west[hi]), BottomTeam: teamFromStandings(west[lo]),
			}
		} else {
			westR1[i] = tbdSeries(1, "West")
		}
	}

	series := make([]types.PlayoffSeries, 15)
	copy(series[0:4], eastR1)
	series[4] = tbdSeries(2, "East")
	series[5] = tbdSeries(2, "East")
	series[6] = tbdSeries(3, "East")
	series[7] = tbdSeries(4, "Finals")
	series[8] = tbdSeries(3, "West")
	series[9] = tbdSeries(2, "West")
	series[10] = tbdSeries(2, "West")
	copy(series[11:15], westR1)

	return types.PlayoffBracket{Season: season, Series: series}
}

// nbaTeamTricodes maps stable NBA team IDs to official 2-3 letter tricodes.
// The standings API does not include a tricode field; TeamIDs are stable across seasons.
var nbaTeamTricodes = map[int]string{
	1610612737: "ATL", 1610612738: "BOS", 1610612739: "CLE", 1610612740: "NOP",
	1610612741: "CHI", 1610612742: "DAL", 1610612743: "DEN", 1610612744: "GSW",
	1610612745: "HOU", 1610612746: "LAC", 1610612747: "LAL", 1610612748: "MIA",
	1610612749: "MIL", 1610612750: "MIN", 1610612751: "BKN", 1610612752: "NYK",
	1610612753: "ORL", 1610612754: "IND", 1610612755: "PHI", 1610612756: "PHX",
	1610612757: "POR", 1610612758: "SAC", 1610612759: "SAS", 1610612760: "OKC",
	1610612761: "TOR", 1610612762: "UTA", 1610612763: "MEM", 1610612764: "WAS",
	1610612765: "DET", 1610612766: "CHA",
}

func teamTricode(t types.Team) string {
	if tri, ok := nbaTeamTricodes[t.TeamID]; ok {
		return tri
	}
	if len(t.TeamName) >= 3 {
		return t.TeamName[:3]
	}
	return t.TeamName
}

// PopulatePlayoffBracket converts a commonPlayoffSeries API response into a PlayoffBracket.
// It derives round/conference from the SERIES_ID format and determines series winners by
// cross-referencing which teams appear in the next round.
func PopulatePlayoffBracket(rs types.ResponseSet, season string) (types.PlayoffBracket, error) {
	if len(rs.ResultSets) == 0 || rs.ResultSets[0].Headers == nil {
		return types.PlayoffBracket{}, fmt.Errorf("empty resultSet for playoff bracket")
	}
	rs0 := rs.ResultSets[0]
	if len(rs0.RowSet) == 0 {
		return types.PlayoffBracket{}, fmt.Errorf("no playoff series data for season %s", season)
	}

	headerIdx := make(map[string]int, len(rs0.Headers))
	for i, h := range rs0.Headers {
		headerIdx[h] = i
	}
	for _, req := range []string{"HOME_TEAM_ID", "VISITOR_TEAM_ID", "SERIES_ID", "GAME_NUM"} {
		if _, ok := headerIdx[req]; !ok {
			return types.PlayoffBracket{}, fmt.Errorf("missing required header: %s", req)
		}
	}

	type seriesEntry struct {
		homeID  string
		visitID string
		maxGame int
	}

	seriesMap := make(map[string]*seriesEntry)
	for _, row := range rs0.RowSet {
		sid := rowString(row[headerIdx["SERIES_ID"]])
		gameNum := int(rowInt64(row[headerIdx["GAME_NUM"]]))
		homeID := fmt.Sprintf("%d", rowInt64(row[headerIdx["HOME_TEAM_ID"]]))
		visitID := fmt.Sprintf("%d", rowInt64(row[headerIdx["VISITOR_TEAM_ID"]]))

		if _, exists := seriesMap[sid]; !exists {
			seriesMap[sid] = &seriesEntry{}
		}
		se := seriesMap[sid]
		if gameNum == 1 {
			se.homeID = homeID
			se.visitID = visitID
		}
		if gameNum > se.maxGame {
			se.maxGame = gameNum
		}
	}

	// Collect team IDs per round to cross-reference winners.
	roundTeams := [5]map[string]bool{}
	for i := range roundTeams {
		roundTeams[i] = make(map[string]bool)
	}
	for sid, se := range seriesMap {
		if len(sid) < 9 {
			continue
		}
		round := int(sid[7] - '0')
		if round >= 1 && round <= 4 {
			roundTeams[round][se.homeID] = true
			roundTeams[round][se.visitID] = true
		}
	}

	tbd := types.PlayoffTeam{Tricode: "TBD"}
	bracket := make([]types.PlayoffSeries, 15)
	for i := range bracket {
		bracket[i] = types.PlayoffSeries{Status: "pre", TopTeam: tbd, BottomTeam: tbd}
	}

	for sid, se := range seriesMap {
		if len(sid) < 9 {
			continue
		}
		round := int(sid[7] - '0')
		seriesIdx := int(sid[8] - '0')
		if round < 1 || round > 4 {
			continue
		}

		cursorIdx := playoffCursorIndex(round, seriesIdx)
		if cursorIdx < 0 || cursorIdx >= 15 {
			continue
		}

		top := types.PlayoffTeam{TeamID: se.homeID, Tricode: nbaTeamTricodeByID(se.homeID)}
		bot := types.PlayoffTeam{TeamID: se.visitID, Tricode: nbaTeamTricodeByID(se.visitID)}
		status := playoffSeriesStatus(round, se.maxGame, top.TeamID, bot.TeamID, roundTeams)

		if status == "complete" && round < 4 {
			switch {
			case roundTeams[round+1][top.TeamID]:
				top.Wins = 4
				bot.Wins = se.maxGame - 4
			case roundTeams[round+1][bot.TeamID]:
				bot.Wins = 4
				top.Wins = se.maxGame - 4
			}
		}

		bracket[cursorIdx] = types.PlayoffSeries{
			SeriesID:   sid,
			Round:      round,
			Conference: playoffConference(round, seriesIdx),
			TopTeam:    top,
			BottomTeam: bot,
			Status:     status,
		}
	}

	return types.PlayoffBracket{Season: season, Series: bracket}, nil
}

func playoffSeriesStatus(round, maxGame int, topID, botID string, roundTeams [5]map[string]bool) string {
	if maxGame == 0 {
		return "pre"
	}
	if round < 4 && (roundTeams[round+1][topID] || roundTeams[round+1][botID]) {
		return "complete"
	}
	if maxGame >= 4 {
		return "complete"
	}
	return "active"
}

// playoffCursorIndex maps (round, NBA series index) to the bracket cursor index.
// Round 1: East series 0-3, West series 4-7. Within each conference: 0=1v8, 1=2v7, 2=3v6, 3=4v5.
// Bracket cursor order within East R1: 0=1v8, 1=4v5, 2=3v6, 3=2v7.
func playoffCursorIndex(round, seriesIdx int) int {
	switch round {
	case 1:
		eastCursor := [4]int{0, 3, 2, 1}
		westCursor := [4]int{11, 14, 13, 12}
		if seriesIdx <= 3 {
			return eastCursor[seriesIdx]
		}
		if seriesIdx <= 7 {
			return westCursor[seriesIdx-4]
		}
	case 2:
		switch seriesIdx {
		case 0:
			return 4
		case 1:
			return 5
		case 2:
			return 9
		case 3:
			return 10
		}
	case 3:
		if seriesIdx == 0 {
			return 6
		}
		return 8
	case 4:
		return 7
	}
	return -1
}

func playoffConference(round, seriesIdx int) string {
	switch round {
	case 1:
		if seriesIdx <= 3 {
			return "East"
		}
		return "West"
	case 2:
		if seriesIdx <= 1 {
			return "East"
		}
		return "West"
	case 3:
		if seriesIdx == 0 {
			return "East"
		}
		return "West"
	}
	return "Finals"
}

func nbaTeamTricodeByID(teamID string) string {
	id, err := strconv.Atoi(teamID)
	if err != nil {
		return ""
	}
	return nbaTeamTricodes[id]
}

func rowInt64(v interface{}) int64 {
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int64:
		return n
	case int:
		return int64(n)
	}
	return 0
}

// rowString extracts a string value from an interface{} row cell.
// It handles both native string values and numeric values (returned as zero-padded strings).
func rowString(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	case float64:
		return fmt.Sprintf("%09d", int64(s))
	case int64:
		return fmt.Sprintf("%09d", s)
	case int:
		return fmt.Sprintf("%09d", int64(s))
	}
	return ""
}

// PopulatePlayoffSeriesGames converts a commonPlayoffSeries API response into game records.
// seriesID filters results to a single series; pass "" to return all games.
func PopulatePlayoffSeriesGames(rs types.ResponseSet, seriesID string) ([]types.PlayoffGame, error) {
	if len(rs.ResultSets) == 0 || rs.ResultSets[0].Headers == nil {
		return nil, fmt.Errorf("empty resultSet for playoff series games")
	}
	// TODO: implement field extraction after inspecting a real commonPlayoffSeries response.
	return nil, nil
}
