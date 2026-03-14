package tui

import "testing"

func TestSeasonSelector_PrevSeason(t *testing.T) {
	ss := NewSeasonSelector("2024-25")
	ss.prevSeason()
	if ss.season != "2023-24" {
		t.Errorf("prevSeason() = %s, want 2023-24", ss.season)
	}
}

func TestSeasonSelector_PrevSeason_AtFloor(t *testing.T) {
	ss := NewSeasonSelector("2000-01")
	ss.prevSeason()
	if ss.season != "2000-01" {
		t.Errorf("prevSeason() at floor = %s, want 2000-01", ss.season)
	}
}

func TestSeasonSelector_NextSeason_AtCeiling(t *testing.T) {
	ss := NewSeasonSelector("2025-26")
	ss.ceiling = "2025-26"
	ss.nextSeason()
	if ss.season != "2025-26" {
		t.Errorf("nextSeason() at ceiling = %s, want 2025-26", ss.season)
	}
}

func TestSeasonSelector_FormatSeason(t *testing.T) {
	cases := []struct {
		year int
		want string
	}{
		{2024, "2024-25"},
		{1999, "1999-00"},
		{2099, "2099-00"},
	}
	for _, c := range cases {
		if got := formatSeasonFromYear(c.year); got != c.want {
			t.Errorf("formatSeasonFromYear(%d) = %s, want %s", c.year, got, c.want)
		}
	}
}
