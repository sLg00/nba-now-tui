package datamodels

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func mockUnmarshall(_ string) (ResponseSet, error) {
	return ResponseSet{
		ResultSets: []ResultSet{
			{
				Headers: []string{
					"LeagueID", "SeasonID", "TeamID", "TeamCity", "TeamName", "TeamSlug", "Conference",
					"WINS", "LOSSES",
				},
				RowSet: [][]interface{}{
					{"00", "2023", 1, "Boston", "Celtics", "boston-celtics", "East", 45, 20},
					{"00", "2023", 2, "Golden State", "Warriors", "golden-state-warriors", "West", 42, 23},
				},
			},
		},
	}, nil
}

func mockUnmarshallError(_ string) (ResponseSet, error) {
	return ResponseSet{
		ResultSets: []ResultSet{
			{
				Headers: nil,
				RowSet:  nil,
			},
		},
	}, fmt.Errorf("mock unmarshall error")
}

func TestPopulateTeamStats(t *testing.T) {
	realArguments := os.Args
	defer func() { os.Args = realArguments }()
	os.Args = []string{"appName", "-d", "2024-12-01"}

	teams, headers, err := PopulateTeamStats(mockUnmarshall)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	//Header validation
	expectedHeaders := []string{"LeagueID", "SeasonID", "TeamID", "TeamCity", "TeamName", "TeamSlug", "Conference", "WINS", "LOSSES"}
	if !reflect.DeepEqual(headers, expectedHeaders) {
		t.Errorf("Expected headers %v, got %v", expectedHeaders, headers)
	}

	// Team count validation
	if len(teams) != 2 {
		t.Fatalf("Expected 2 teams, got %d", len(teams))
	}
}

func TestPopulateTeamStats_Error(t *testing.T) {
	realArguments := os.Args
	defer func() { os.Args = realArguments }()
	os.Args = []string{"appName", "-d", "2024-12-01"}

	teams, headers, err := PopulateTeamStats(mockUnmarshallError)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}

	if teams != nil {
		t.Errorf("Expected nil teams, got: %+v", teams)
	}

	if headers != nil {
		t.Errorf("Expected nil headers, got: %+v", headers)
	}
}

// TestSplitStandingsPerConference checks if teams are split correctly by conference.
func TestSplitStandingsPerConference(t *testing.T) {
	teams := Teams{
		{TeamID: 1, TeamCity: "Boston", TeamName: "Celtics", Conference: "East"},
		{TeamID: 2, TeamCity: "Golden State", TeamName: "Warriors", Conference: "West"},
		{TeamID: 3, TeamCity: "Milwaukee", TeamName: "Bucks", Conference: "East"},
		{TeamID: 4, TeamCity: "Denver", TeamName: "Nuggets", Conference: "West"},
	}

	eastTeams, westTeams := teams.SplitStandingsPerConference()

	// Validate East teams
	if len(eastTeams) != 2 {
		t.Errorf("Expected 2 East teams, got: %d", len(eastTeams))
	}
	if eastTeams[0].TeamName != "Celtics" || eastTeams[1].TeamName != "Bucks" {
		t.Errorf("East teams data mismatch: %+v", eastTeams)
	}

	// Validate West teams
	if len(westTeams) != 2 {
		t.Errorf("Expected 2 West teams, got: %d", len(westTeams))
	}
	if westTeams[0].TeamName != "Warriors" || westTeams[1].TeamName != "Nuggets" {
		t.Errorf("West teams data mismatch: %+v", westTeams)
	}
}

// TestSplitStandingsPerConference_Empty checks behavior when there are no teams.
func TestSplitStandingsPerConference_Empty(t *testing.T) {
	teams := Teams{}

	eastTeams, westTeams := teams.SplitStandingsPerConference()
	if len(eastTeams) != 0 {
		t.Errorf("Expected 0 East teams, got: %d", len(eastTeams))
	}
	if len(westTeams) != 0 {
		t.Errorf("Expected 0 West teams, got: %d", len(westTeams))
	}
}
