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

func TestPopulatePlayoffBracket_ReturnsErrorOnEmptyRowSet(t *testing.T) {
	rs := types.ResponseSet{ResultSets: []types.ResultSet{{
		Headers: []string{"GAME_ID", "HOME_TEAM_ID", "VISITOR_TEAM_ID", "SERIES_ID", "GAME_NUM"},
		RowSet:  [][]interface{}{},
	}}}
	_, err := PopulatePlayoffBracket(rs, "2025-26")
	if err == nil {
		t.Error("expected error for empty rowSet (no playoff data)")
	}
}

func TestPopulatePlayoffBracket_BuildsBracket(t *testing.T) {
	// East R1 series 0 (1v8): BOS(1610612738) home vs MIA(1610612748), 4 games → BOS sweeps
	// East R2 series 0: BOS home vs OKC — proves BOS won R1 (cross-reference)
	makeRow := func(homeID, visitID, seriesID, gameNum float64) []interface{} {
		return []interface{}{float64(0), homeID, visitID, seriesID, gameNum}
	}
	rows := [][]interface{}{
		makeRow(1610612738, 1610612748, 4230010, 1), // R1 East 1v8 G1
		makeRow(1610612738, 1610612748, 4230010, 2), // G2
		makeRow(1610612748, 1610612738, 4230010, 3), // G3
		makeRow(1610612748, 1610612738, 4230010, 4), // G4
		makeRow(1610612738, 1610612760, 4230020, 1), // R2 East: BOS in R2 → proves BOS won R1
	}
	rs := types.ResponseSet{ResultSets: []types.ResultSet{{
		Headers: []string{"GAME_ID", "HOME_TEAM_ID", "VISITOR_TEAM_ID", "SERIES_ID", "GAME_NUM"},
		RowSet:  rows,
	}}}

	bracket, err := PopulatePlayoffBracket(rs, "2023-24")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bracket.Season != "2023-24" {
		t.Errorf("Season = %s, want 2023-24", bracket.Season)
	}

	// cursor 0 = East R1 series 0 (1v8)
	s := bracket.Series[0]
	if s.TopTeam.Tricode != "BOS" {
		t.Errorf("East R1 TopTeam = %s, want BOS", s.TopTeam.Tricode)
	}
	if s.BottomTeam.Tricode != "MIA" {
		t.Errorf("East R1 BottomTeam = %s, want MIA", s.BottomTeam.Tricode)
	}
	if s.TopTeam.Wins != 4 {
		t.Errorf("BOS wins = %d, want 4 (advanced to R2)", s.TopTeam.Wins)
	}
	if s.BottomTeam.Wins != 0 {
		t.Errorf("MIA wins = %d, want 0 (4 games - 4)", s.BottomTeam.Wins)
	}
	if s.Status != "complete" {
		t.Errorf("R1 status = %s, want complete", s.Status)
	}

	// cursor 4 = East R2 series 0 (upper semi)
	semi := bracket.Series[4]
	if semi.TopTeam.Tricode != "BOS" {
		t.Errorf("East Semi TopTeam = %s, want BOS", semi.TopTeam.Tricode)
	}
}

// TestPopulatePlayoffBracket_BuildsBracket_StringSeriesID verifies that SERIES_ID values
// returned as strings by the real NBA API (e.g. "004230010") are parsed correctly.
func TestPopulatePlayoffBracket_BuildsBracket_StringSeriesID(t *testing.T) {
	makeRow := func(homeID, visitID float64, seriesID string, gameNum float64) []interface{} {
		return []interface{}{float64(0), homeID, visitID, seriesID, gameNum}
	}
	rows := [][]interface{}{
		makeRow(1610612738, 1610612748, "004230010", 1), // R1 East 1v8 G1
		makeRow(1610612738, 1610612748, "004230010", 2), // G2
		makeRow(1610612748, 1610612738, "004230010", 3), // G3
		makeRow(1610612748, 1610612738, "004230010", 4), // G4
		makeRow(1610612738, 1610612760, "004230020", 1), // R2 East: BOS in R2
	}
	rs := types.ResponseSet{ResultSets: []types.ResultSet{{
		Headers: []string{"GAME_ID", "HOME_TEAM_ID", "VISITOR_TEAM_ID", "SERIES_ID", "GAME_NUM"},
		RowSet:  rows,
	}}}

	bracket, err := PopulatePlayoffBracket(rs, "2023-24")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := bracket.Series[0]
	if s.TopTeam.Tricode != "BOS" {
		t.Errorf("East R1 TopTeam = %s, want BOS", s.TopTeam.Tricode)
	}
	if s.BottomTeam.Tricode != "MIA" {
		t.Errorf("East R1 BottomTeam = %s, want MIA", s.BottomTeam.Tricode)
	}
	if s.Status != "complete" {
		t.Errorf("R1 status = %s, want complete", s.Status)
	}
	semi := bracket.Series[4]
	if semi.TopTeam.Tricode != "BOS" {
		t.Errorf("East Semi TopTeam = %s, want BOS", semi.TopTeam.Tricode)
	}
}

func TestPopulatePlayoffSeriesGames_ReturnsErrorOnEmptyResultSet(t *testing.T) {
	rs := types.ResponseSet{ResultSets: []types.ResultSet{{Headers: nil, RowSet: nil}}}
	_, err := PopulatePlayoffSeriesGames(rs, "")
	if err == nil {
		t.Error("expected error for empty ResultSet")
	}
}
