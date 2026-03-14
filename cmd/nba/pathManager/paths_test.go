package pathManager

import (
	"strings"
	"testing"
)

type mockDateProvider struct {
	date   string
	season string
}

func (m *mockDateProvider) GetCurrentDate() (string, error) { return m.date, nil }
func (m *mockDateProvider) GetCurrentSeason() string        { return m.season }

func TestGetFullPath_PlayerInfo(t *testing.T) {
	pm := PathFactory(&mockDateProvider{date: "2025-02-14", season: "2024-25"}, "")
	got := pm.GetFullPath("playerInfo", "1628389")
	if !strings.HasSuffix(got, "/.config/nba-tui/playerprofiles/1628389_info") {
		t.Errorf("GetFullPath(playerInfo) = %s, want suffix playerprofiles/1628389_info", got)
	}
}

func TestGetFullPath_PlayerCareerStats(t *testing.T) {
	pm := PathFactory(&mockDateProvider{date: "2025-02-14", season: "2024-25"}, "")
	got := pm.GetFullPath("playerCareerStats", "1628389")
	if !strings.HasSuffix(got, "/.config/nba-tui/playerprofiles/1628389_career") {
		t.Errorf("GetFullPath(playerCareerStats) = %s, want suffix playerprofiles/1628389_career", got)
	}
}

func TestGetFullPath_PlayerGameLog(t *testing.T) {
	pm := PathFactory(&mockDateProvider{date: "2025-02-14", season: "2024-25"}, "")
	got := pm.GetFullPath("playerGameLog", "1628389")
	if !strings.HasSuffix(got, "/.config/nba-tui/playerprofiles/1628389_gamelog") {
		t.Errorf("GetFullPath(playerGameLog) = %s, want suffix playerprofiles/1628389_gamelog", got)
	}
}

func TestGetBasePaths_IncludesPlayerProfiles(t *testing.T) {
	pm := PathFactory(&mockDateProvider{date: "2025-02-14", season: "2024-25"}, "")
	paths := pm.GetBasePaths()
	found := false
	for _, p := range paths {
		if strings.HasSuffix(p, "playerprofiles/") {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetBasePaths() should include playerprofiles/ path")
	}
}

func TestGetFullPath_PlayoffBracket(t *testing.T) {
	p := &PathComps{
		Home:         "/home/user",
		Path:         "/.config/nba-tui/",
		PlayoffsPath: "playoffs/",
	}
	got := p.GetFullPath("playoffBracket", "2023-24")
	want := "/home/user/.config/nba-tui/playoffs/2023-24_bracket"
	if got != want {
		t.Errorf("GetFullPath(playoffBracket) = %s, want %s", got, want)
	}
}

func TestGetFullPath_PlayoffSeriesGames(t *testing.T) {
	p := &PathComps{
		Home:         "/home/user",
		Path:         "/.config/nba-tui/",
		PlayoffsPath: "playoffs/",
	}
	got := p.GetFullPath("playoffSeriesGames", "0042300401")
	want := "/home/user/.config/nba-tui/playoffs/0042300401_games"
	if got != want {
		t.Errorf("GetFullPath(playoffSeriesGames) = %s, want %s", got, want)
	}
}

func TestGetBasePaths_IncludesPlayoffs(t *testing.T) {
	p := &PathComps{
		Home:         "/home/user",
		Path:         "/.config/nba-tui/",
		PlayoffsPath: "playoffs/",
	}
	paths := p.GetBasePaths()
	for _, path := range paths {
		if path == "/home/user/.config/nba-tui/playoffs/" {
			return
		}
	}
	t.Error("GetBasePaths() does not include playoffs path")
}
