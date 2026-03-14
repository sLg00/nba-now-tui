package converters

import (
	"testing"

	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

func TestProjectedBracketFromStandings_CorrectMatchups(t *testing.T) {
	east := make([]types.Team, 8)
	west := make([]types.Team, 8)
	for i := 0; i < 8; i++ {
		east[i] = types.Team{TeamCity: "ECity", TeamName: "ETeam", PlayoffSeeding: i + 1, Conference: "East"}
		west[i] = types.Team{TeamCity: "WCity", TeamName: "WTeam", PlayoffSeeding: i + 1, Conference: "West"}
	}

	bracket := ProjectedBracketFromStandings(east, west, "2024-25")

	if bracket.Season != "2024-25" {
		t.Errorf("Season = %s, want 2024-25", bracket.Season)
	}

	if len(bracket.Series) != 15 {
		t.Errorf("len(Series) = %d, want 15", len(bracket.Series))
	}

	// East R1 index 0: seed 1 vs seed 8
	e1 := bracket.Series[0]
	if e1.TopTeam.Seed != 1 || e1.BottomTeam.Seed != 8 {
		t.Errorf("East R1 series 0: seeds %d vs %d, want 1 vs 8", e1.TopTeam.Seed, e1.BottomTeam.Seed)
	}
	if e1.Status != "pre" {
		t.Errorf("East R1 series 0 status = %s, want pre", e1.Status)
	}

	// East R1 index 1: seed 4 vs seed 5
	e2 := bracket.Series[1]
	if e2.TopTeam.Seed != 4 || e2.BottomTeam.Seed != 5 {
		t.Errorf("East R1 series 1: seeds %d vs %d, want 4 vs 5", e2.TopTeam.Seed, e2.BottomTeam.Seed)
	}

	// Finals (index 7) should be TBD
	finals := bracket.Series[7]
	if finals.TopTeam.Tricode != "TBD" {
		t.Errorf("Finals TopTeam = %s, want TBD", finals.TopTeam.Tricode)
	}
}

func TestPopulatePlayoffBracket_ReturnsErrorOnEmptyResultSet(t *testing.T) {
	rs := types.ResponseSet{ResultSets: []types.ResultSet{{Headers: nil, RowSet: nil}}}
	_, err := PopulatePlayoffBracket(rs, "2023-24")
	if err == nil {
		t.Error("expected error for empty ResultSet")
	}
}

func TestPopulatePlayoffSeriesGames_ReturnsErrorOnEmptyResultSet(t *testing.T) {
	rs := types.ResponseSet{ResultSets: []types.ResultSet{{Headers: nil, RowSet: nil}}}
	_, err := PopulatePlayoffSeriesGames(rs, "")
	if err == nil {
		t.Error("expected error for empty ResultSet")
	}
}
