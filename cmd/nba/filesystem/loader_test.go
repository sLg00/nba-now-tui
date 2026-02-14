package filesystemops

import (
	"testing"
)

type mockPathManager struct {
	fullPathFunc func(name, param string) string
}

func (m *mockPathManager) GetFullPath(name, param string) string {
	if m.fullPathFunc != nil {
		return m.fullPathFunc(name, param)
	}
	return ""
}

func (m *mockPathManager) GetBasePaths() []string { return nil }

type mockFsHandler struct {
	readFileFunc func(path string) ([]byte, error)
}

func (m *mockFsHandler) ReadFile(path string) ([]byte, error) {
	if m.readFileFunc != nil {
		return m.readFileFunc(path)
	}
	return nil, nil
}

func (m *mockFsHandler) WriteFile(string, []byte) error        { return nil }
func (m *mockFsHandler) FileExists(string) bool                { return false }
func (m *mockFsHandler) EnsureDirectoryExists(string) error    { return nil }
func (m *mockFsHandler) CleanOldFiles([]string) error          { return nil }

func TestLoadPlayerInfo(t *testing.T) {
	json := `{"resultSets":[{"name":"CommonPlayerInfo","headers":["FIRST_NAME"],"rowSet":[["Bam"]]}]}`
	paths := &mockPathManager{
		fullPathFunc: func(name, param string) string {
			if name != "playerInfo" || param != "1628389" {
				t.Errorf("unexpected path args: name=%s, param=%s", name, param)
			}
			return "/tmp/test_playerinfo"
		},
	}
	fs := &mockFsHandler{
		readFileFunc: func(path string) ([]byte, error) {
			return []byte(json), nil
		},
	}

	dl := NewDataLoader(fs, paths)
	rs, err := dl.LoadPlayerInfo("1628389")
	if err != nil {
		t.Fatalf("LoadPlayerInfo() error: %v", err)
	}
	if len(rs.ResultSets) == 0 {
		t.Fatal("LoadPlayerInfo() returned empty ResultSets")
	}
	if rs.ResultSets[0].Name != "CommonPlayerInfo" {
		t.Errorf("expected name CommonPlayerInfo, got %s", rs.ResultSets[0].Name)
	}
}

func TestLoadPlayerCareerStats(t *testing.T) {
	json := `{"resultSets":[{"name":"SeasonTotalsRegularSeason","headers":["SEASON_ID"],"rowSet":[["2024-25"]]}]}`
	paths := &mockPathManager{
		fullPathFunc: func(name, param string) string {
			if name != "playerCareerStats" {
				t.Errorf("unexpected path name: %s", name)
			}
			return "/tmp/test_career"
		},
	}
	fs := &mockFsHandler{
		readFileFunc: func(path string) ([]byte, error) {
			return []byte(json), nil
		},
	}

	dl := NewDataLoader(fs, paths)
	rs, err := dl.LoadPlayerCareerStats("1628389")
	if err != nil {
		t.Fatalf("LoadPlayerCareerStats() error: %v", err)
	}
	if rs.ResultSets[0].Name != "SeasonTotalsRegularSeason" {
		t.Errorf("expected name SeasonTotalsRegularSeason, got %s", rs.ResultSets[0].Name)
	}
}

func TestLoadPlayerGameLog(t *testing.T) {
	json := `{"resultSets":[{"name":"PlayerGameLog","headers":["GAME_DATE"],"rowSet":[["FEB 10, 2025"]]}]}`
	paths := &mockPathManager{
		fullPathFunc: func(name, param string) string {
			if name != "playerGameLog" {
				t.Errorf("unexpected path name: %s", name)
			}
			return "/tmp/test_gamelog"
		},
	}
	fs := &mockFsHandler{
		readFileFunc: func(path string) ([]byte, error) {
			return []byte(json), nil
		},
	}

	dl := NewDataLoader(fs, paths)
	rs, err := dl.LoadPlayerGameLog("1628389")
	if err != nil {
		t.Fatalf("LoadPlayerGameLog() error: %v", err)
	}
	if rs.ResultSets[0].Name != "PlayerGameLog" {
		t.Errorf("expected name PlayerGameLog, got %s", rs.ResultSets[0].Name)
	}
}
