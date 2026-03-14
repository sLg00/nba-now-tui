package types

import "testing"

func TestPlayoffSeries_IsComplete(t *testing.T) {
	s := PlayoffSeries{TopTeam: PlayoffTeam{Wins: 4}, BottomTeam: PlayoffTeam{Wins: 2}}
	if !s.IsComplete() {
		t.Error("expected series with 4 wins to be complete")
	}
	s2 := PlayoffSeries{TopTeam: PlayoffTeam{Wins: 2}, BottomTeam: PlayoffTeam{Wins: 3}}
	if s2.IsComplete() {
		t.Error("expected series with max 3 wins to be incomplete")
	}
}

func TestPlayoffSeries_Leader(t *testing.T) {
	s := PlayoffSeries{
		TopTeam:    PlayoffTeam{Tricode: "BOS", Wins: 3},
		BottomTeam: PlayoffTeam{Tricode: "MIA", Wins: 1},
	}
	if got := s.Leader(); got != "BOS" {
		t.Errorf("Leader() = %s, want BOS", got)
	}
	tied := PlayoffSeries{
		TopTeam:    PlayoffTeam{Tricode: "BOS", Wins: 2},
		BottomTeam: PlayoffTeam{Tricode: "MIA", Wins: 2},
	}
	if got := tied.Leader(); got != "" {
		t.Errorf("Leader() tied series = %s, want empty string", got)
	}
}
