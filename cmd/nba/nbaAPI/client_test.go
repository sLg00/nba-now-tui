package nbaAPI

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type MockDateProvider struct {
	currentDate   string
	currentSeason string
	dateError     error
}

func (m *MockDateProvider) GetCurrentDate() (string, error) {
	return m.currentDate, m.dateError
}

func (m *MockDateProvider) GetCurrentSeason() string {
	return m.currentSeason
}

type MockHTTPClient struct {
	getFunc        func(url RequestURL) ([]byte, error)
	setHeadersFunc func() http.Header
}

func (m *MockHTTPClient) Get(url RequestURL) ([]byte, error) {
	if m.getFunc != nil {
		return m.getFunc(url)
	}
	return nil, nil
}

func (m *MockHTTPClient) SetHeaders() http.Header {
	if m.setHeadersFunc != nil {
		return m.setHeadersFunc()
	}
	return http.Header{}
}

type MockFileSystem struct {
	writeFileFunc     func(path string, data []byte) error
	readFileFunc      func(path string) ([]byte, error)
	fileExistsFunc    func(path string) bool
	cleanOldFilesFunc func(path []string) error
}

func (m *MockFileSystem) writeFile(path string, data []byte) error {
	if m.writeFileFunc != nil {
		return m.writeFileFunc(path, data)
	}
	return nil
}

func (m *MockFileSystem) readFile(path string) ([]byte, error) {
	if m.readFileFunc != nil {
		return m.readFileFunc(path)
	}
	return nil, nil
}

func (m *MockFileSystem) fileExists(path string) bool {
	if m.fileExistsFunc != nil {
		return m.fileExistsFunc(path)
	}
	return false
}

func (m *MockFileSystem) cleanOldFiles(path []string) error {
	if m.cleanOldFilesFunc != nil {
		return m.cleanOldFilesFunc(path)
	}
	return nil
}

type MockPathManager struct {
	basePathsFunc func() []string
	fullPathFunc  func(name, param string) string
}

func (m *MockPathManager) GetBasePaths() []string {
	if m.basePathsFunc != nil {
		return m.basePathsFunc()
	}
	return []string{}
}

func (m *MockPathManager) GetFullPath(name, param string) string {
	if m.fullPathFunc != nil {
		return m.fullPathFunc(name, param)
	}
	return ""
}

func TestNbaRequestBuilder_BuildLeagueLeadersRequest(t *testing.T) {
	mockDates := &MockDateProvider{
		currentSeason: "2024-25",
	}
	baseUrl := BaseURL
	rb := NewRequestBuilder(baseUrl, mockDates)

	got := string(rb.BuildLeagueLeadersRequest())
	want := "https://stats.nba.com/stats/leagueleaders" +
		"?LeagueID=00&PerMode=PerGame&Scope=S&Season=2024-25" +
		"&SeasonType=Regular+Season&StatCategory=PTS"

	if !urlsEqual(t, want, got) {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestNbaRequestBuilder_BuildDailyScoresRequest(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		dateErr error
		want    string
		wantErr error
	}{
		{
			name: "successful request",
			date: "2025-02-20",
			want: "https://stats.nba.com/stats/scoreboardv2" +
				"?DayOffset=0&GameDate=2025-02-20&LeagueID=00",
		},
		{
			name:    "date error",
			date:    "",
			dateErr: errors.New("date error"),
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDates := &MockDateProvider{
				currentDate: tt.date,
				dateError:   tt.dateErr}

			rb := NewRequestBuilder(BaseURL, mockDates)
			got := string(rb.BuildDailyScoresRequest())

			if !urlsEqual(t, tt.want, got) {
				t.Errorf("BildDailyscoresRequest() got %s, wanted %s", got, tt.want)
			}
		})
	}
}

func TestNbaRequestBuilder_BuildBoxScoreRequest(t *testing.T) {
	tests := []struct {
		name      string
		want      string
		gameID    string
		gameIDErr error
		date      string
	}{
		{
			name:   "successful request",
			gameID: "0052300101",
			date:   "",
			want: "https://stats.nba.com/stats/boxscoretraditionalv3?" +
				"EndPeriod=4&EndRange=0&GameID=0052300101&RangeType=0&StartPeriod=1&StartRange=0",
		},
		{
			name:      "gameID not found",
			gameID:    "",
			want:      "",
			gameIDErr: errors.New("gameID not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDates := &MockDateProvider{
				currentDate: tt.date,
			}
			rb := NewRequestBuilder(BaseURL, mockDates)
			got := string(rb.BuildBoxScoreRequest(tt.gameID))
			if !urlsEqual(t, tt.want, got) {
				t.Errorf("BuildBoxScoreRequest() got %s, wanted %s", got, tt.want)
			}
		})
	}
}

func TestClient_MakeDefaultRequests(t *testing.T) {
	return
}

func TestClient_FetchBoxScore(t *testing.T) {
	return
}

// urlsEqual is a helper function for URL comparison tests
func urlsEqual(t *testing.T, got, want string) bool {
	if got == want {
		return true
	}

	gotURL, err1 := url.Parse(got)
	wantURL, err2 := url.Parse(want)

	if err1 != nil || err2 != nil {
		t.Errorf("failed to parse urls: got=%v, want=%v", err1, err2)
		return false
	}

	if gotURL.Scheme != wantURL.Scheme ||
		gotURL.Host != wantURL.Host ||
		gotURL.Path != wantURL.Path {
		return false
	}

	gotParams := gotURL.Query()
	wantParams := wantURL.Query()
	return reflect.DeepEqual(gotParams, wantParams)
}
