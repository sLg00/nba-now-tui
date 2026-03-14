package converters

import (
	"fmt"
	"sort"

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

// PopulatePlayoffBracket converts a leagueSeriesStandings API response into a PlayoffBracket.
// Field extraction is deferred until a real API response can be inspected at runtime.
func PopulatePlayoffBracket(rs types.ResponseSet, season string) (types.PlayoffBracket, error) {
	if len(rs.ResultSets) == 0 || rs.ResultSets[0].Headers == nil {
		return types.PlayoffBracket{}, fmt.Errorf("empty resultSet for playoff bracket")
	}
	// TODO: implement field extraction after inspecting a real leagueSeriesStandings response.
	// Follow the header-index lookup pattern used in PopulateTeamStats (teamDataConverters.go).
	return types.PlayoffBracket{Season: season}, nil
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
