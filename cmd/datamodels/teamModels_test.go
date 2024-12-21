package datamodels

import (
	"testing"
)

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
