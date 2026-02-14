package nbaAPI

import (
	"errors"
	"fmt"
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

type MockRequestBuilder struct {
	buildRequests                   func(param string) map[string]RequestURL
	buildLeagueLeadersRequests      func() RequestURL
	buildSeasonStandingsRequests    func() RequestURL
	buildDailyScoresRequests        func() RequestURL
	buildDailyScoresForDateRequests func(date string) RequestURL
	buildBoxScoreRequests           func(gameID string) RequestURL
	buildTeamInfoRequests           func(teamID string) RequestURL
	buildPlayerIndexRequests        func(teamID string) RequestURL
	buildPlayerInfoRequests         func(playerID string) RequestURL
	buildPlayerCareerStatsRequests  func(playerID string) RequestURL
	buildPlayerGameLogRequests      func(playerID string) RequestURL
}

type MockHTTPClient struct {
	getFunc        func(url RequestURL) ([]byte, error)
	setHeadersFunc func() http.Header
}

type MockFileSystem struct {
	writeFileFunc     func(path string, data []byte) error
	readFileFunc      func(path string) ([]byte, error)
	fileExistsFunc    func(path string) bool
	cleanOldFilesFunc func(path []string) error
	dirExistsFunc     func(path string) error
}

type MockPathManager struct {
	basePathsFunc func() []string
	fullPathFunc  func(name, param string) string
}

func (m *MockRequestBuilder) BuildRequests(param string) map[string]RequestURL {
	if m.buildRequests != nil {
		return m.buildRequests(param)
	}
	return map[string]RequestURL{
		"dailyScores":     "https://example.com/scores",
		"leagueLeaders":   "https://example.com/leaders",
		"seasonStandings": "https://example.com/standings",
	}
}

func (m *MockRequestBuilder) BuildLeagueLeadersRequest() RequestURL {
	if m.buildLeagueLeadersRequests != nil {
		return m.buildLeagueLeadersRequests()
	}
	return "https://example.com/leaders"
}

func (m *MockRequestBuilder) BuildSeasonStandingsRequest() RequestURL {
	if m.buildSeasonStandingsRequests != nil {
		return m.buildSeasonStandingsRequests()
	}
	return "https://example.com/standings"
}

func (m *MockRequestBuilder) BuildDailyScoresRequest() RequestURL {
	if m.buildDailyScoresRequests != nil {
		return m.buildDailyScoresRequests()
	}
	return "https://example.com/scores"
}

func (m *MockRequestBuilder) BuildDailyScoresRequestForDate(date string) RequestURL {
	if m.buildDailyScoresForDateRequests != nil {
		return m.buildDailyScoresForDateRequests(date)
	}
	return RequestURL("https://example.com/scores?GameDate=" + date)
}

func (m *MockRequestBuilder) BuildBoxScoreRequest(gameID string) RequestURL {
	if m.buildBoxScoreRequests != nil {
		return m.buildBoxScoreRequests(gameID)
	}
	return "https://example.com/boxscore"
}

func (m *MockRequestBuilder) BuildTeamInfoRequest(teamID string) RequestURL {
	if m.buildTeamInfoRequests != nil {
		return m.buildTeamInfoRequests(teamID)
	}
	return "https://example.com/teaminfo"
}

func (m *MockRequestBuilder) BuildPlayerIndexRequest(teamID string) RequestURL {
	if m.buildPlayerIndexRequests != nil {
		return m.buildPlayerIndexRequests(teamID)
	}
	return "https://example.com/playerindex"
}

func (m *MockRequestBuilder) BuildPlayerInfoRequest(playerID string) RequestURL {
	if m.buildPlayerInfoRequests != nil {
		return m.buildPlayerInfoRequests(playerID)
	}
	return RequestURL("https://example.com/playerinfo?PlayerID=" + playerID)
}

func (m *MockRequestBuilder) BuildPlayerCareerStatsRequest(playerID string) RequestURL {
	if m.buildPlayerCareerStatsRequests != nil {
		return m.buildPlayerCareerStatsRequests(playerID)
	}
	return RequestURL("https://example.com/playercareerstats?PlayerID=" + playerID)
}

func (m *MockRequestBuilder) BuildPlayerGameLogRequest(playerID string) RequestURL {
	if m.buildPlayerGameLogRequests != nil {
		return m.buildPlayerGameLogRequests(playerID)
	}
	return RequestURL("https://example.com/playergamelog?PlayerID=" + playerID)
}

func (m *MockDateProvider) GetCurrentDate() (string, error) {
	return m.currentDate, m.dateError
}

func (m *MockDateProvider) GetCurrentSeason() string {
	return m.currentSeason
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

func (m *MockFileSystem) WriteFile(path string, data []byte) error {
	if m.writeFileFunc != nil {
		return m.writeFileFunc(path, data)
	}
	return nil
}

func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
	if m.readFileFunc != nil {
		return m.readFileFunc(path)
	}
	return nil, nil
}

func (m *MockFileSystem) FileExists(path string) bool {
	if m.fileExistsFunc != nil {
		return m.fileExistsFunc(path)
	}
	return false
}

func (m *MockFileSystem) EnsureDirectoryExists(dir string) error {
	if m.dirExistsFunc != nil {
		return m.dirExistsFunc(dir)
	}
	return nil
}

func (m *MockFileSystem) CleanOldFiles(path []string) error {
	if m.cleanOldFilesFunc != nil {
		return m.cleanOldFilesFunc(path)
	}
	return nil
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

func TestClient_MakeRequests(t *testing.T) {
	tests := []struct {
		name         string
		fileExists   bool
		writeFileErr error
		readFileErr  error
		cleanFileErr error
		httpResponse []byte
		httpErr      error
		wantErr      bool
		gameID       string
		teamID       string
	}{
		{
			name:         "successful default request",
			httpResponse: []byte(`{"data":"success"}`),
			httpErr:      nil,
			fileExists:   false,
			wantErr:      false,
		},
		{
			name:       "file exists",
			fileExists: true,
			gameID:     "0052300101",
		},
		{
			name:         "write file error",
			writeFileErr: errors.New("write file error"),
			httpResponse: []byte(`{data}:"success"`),
			fileExists:   false,
			wantErr:      true,
		},
		{
			name:         "http error",
			httpErr:      errors.New("http error"),
			httpResponse: nil,
			fileExists:   false,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTP := &MockHTTPClient{
				getFunc: func(url RequestURL) ([]byte, error) {
					return tt.httpResponse, tt.httpErr
				},
			}

			mockRequestBuilder := &MockRequestBuilder{
				buildRequests: func(param string) map[string]RequestURL {
					return map[string]RequestURL{
						"dailyScores":     "https://stats.nba.com/stats/scoreboardv2",
						"leagueLeaders":   "https://stats.nba.com/stats/leagueleaders",
						"seasonStandings": "https://stats.nba.com/stats/leaguestandingsv3",
					}
				},
			}

			mockFS := &MockFileSystem{
				fileExistsFunc: func(path string) bool {
					return tt.fileExists
				},
				writeFileFunc: func(path string, data []byte) error {
					return tt.writeFileErr
				},
				cleanOldFilesFunc: func(path []string) error {
					return tt.cleanFileErr
				},
				readFileFunc: func(path string) ([]byte, error) {
					return nil, tt.readFileErr
				},
				dirExistsFunc: func(path string) error { return errors.New("dir exists") },
			}

			mockPaths := &MockPathManager{
				basePathsFunc: func() []string {
					return []string{"/tmp/nba"}
				},
				fullPathFunc: func(name, param string) string {
					return fmt.Sprintf("/tmp/nba/%s", name)
				},
			}

			mockDates := &MockDateProvider{
				currentDate:   "2025-01-02",
				currentSeason: "2024-25",
			}

			client := &Client{
				http:       mockHTTP,
				requests:   mockRequestBuilder,
				Dates:      mockDates,
				Paths:      mockPaths,
				FileSystem: mockFS,
			}
			err := client.MakeDefaultRequests()
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeDefaultRequests() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_FetchBoxScore(t *testing.T) {
	tests := []struct {
		name         string
		httpResponse []byte
		httpErr      error
		writeFileErr error
		wantErr      bool
		fileExists   bool
		gameID       string
	}{
		{
			name:         "successful bs request",
			httpResponse: []byte(`{"data":"success"}`),
			httpErr:      nil,
			wantErr:      false,
			fileExists:   false,
			writeFileErr: nil,
			gameID:       "0052300101",
		},
		{
			name:         "failed bs request",
			httpResponse: nil,
			fileExists:   false,
			writeFileErr: nil,
			httpErr:      errors.New("http error"),
			wantErr:      true,
			gameID:       "4378423",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTP := &MockHTTPClient{
				getFunc: func(url RequestURL) ([]byte, error) {
					return tt.httpResponse, tt.httpErr
				},
			}

			mockRequestBuilder := &MockRequestBuilder{
				buildRequests: func(param string) map[string]RequestURL {
					return map[string]RequestURL{
						"boxScore": RequestURL("https://stats.nba.com/stats/boxscoretraditionalv3?EndPeriod=4&EndRange=0&GameID=" + tt.gameID + "&RangeType=0&StartPeriod=1&StartRange=0"),
					}
				},
			}

			mockPaths := &MockPathManager{
				basePathsFunc: func() []string {
					return []string{"/tmp/nba"}
				},
				fullPathFunc: func(name, param string) string {
					return fmt.Sprintf("/tmp/nba/%s", name)
				},
			}

			mockDates := &MockDateProvider{
				currentDate:   "2025-01-02",
				currentSeason: "2024-25",
			}

			mockFS := &MockFileSystem{
				fileExistsFunc: func(path string) bool {
					return tt.fileExists
				},
				writeFileFunc: func(path string, data []byte) error {
					return tt.writeFileErr
				},
				dirExistsFunc: func(path string) error { return errors.New("dir exists") },
			}

			client := &Client{
				http:       mockHTTP,
				requests:   mockRequestBuilder,
				Dates:      mockDates,
				Paths:      mockPaths,
				FileSystem: mockFS,
			}

			err := client.FetchBoxScore(tt.gameID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchBoxScore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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

func TestClient_FetchDailyScoresForDate(t *testing.T) {
	tests := []struct {
		name         string
		date         string
		httpResponse []byte
		httpErr      error
		writeFileErr error
		wantErr      bool
	}{
		{
			name:         "successful fetch for specific date",
			date:         "2025-01-15",
			httpResponse: []byte(`{"data":"success"}`),
			wantErr:      false,
		},
		{
			name:    "http error",
			date:    "2025-01-15",
			httpErr: errors.New("http error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTP := &MockHTTPClient{
				getFunc: func(url RequestURL) ([]byte, error) {
					return tt.httpResponse, tt.httpErr
				},
			}

			mockRequestBuilder := &MockRequestBuilder{}

			mockFS := &MockFileSystem{
				writeFileFunc: func(path string, data []byte) error {
					return tt.writeFileErr
				},
			}

			client := &Client{
				http:       mockHTTP,
				requests:   mockRequestBuilder,
				FileSystem: mockFS,
			}

			err := client.FetchDailyScoresForDate(tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchDailyScoresForDate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNbaRequestBuilder_BuildPlayerInfoRequest(t *testing.T) {
	mockDates := &MockDateProvider{currentSeason: "2024-25"}
	rb := NewRequestBuilder(BaseURL, mockDates)

	got := string(rb.BuildPlayerInfoRequest("1628389"))
	want := "https://stats.nba.com/stats/commonplayerinfo?PlayerID=1628389"

	if !urlsEqual(t, want, got) {
		t.Errorf("BuildPlayerInfoRequest() got %s, want %s", got, want)
	}
}

func TestNbaRequestBuilder_BuildPlayerCareerStatsRequest(t *testing.T) {
	mockDates := &MockDateProvider{currentSeason: "2024-25"}
	rb := NewRequestBuilder(BaseURL, mockDates)

	got := string(rb.BuildPlayerCareerStatsRequest("1628389"))
	want := "https://stats.nba.com/stats/playercareerstats?PerMode=PerGame&PlayerID=1628389"

	if !urlsEqual(t, want, got) {
		t.Errorf("BuildPlayerCareerStatsRequest() got %s, want %s", got, want)
	}
}

func TestNbaRequestBuilder_BuildPlayerGameLogRequest(t *testing.T) {
	mockDates := &MockDateProvider{currentSeason: "2024-25"}
	rb := NewRequestBuilder(BaseURL, mockDates)

	got := string(rb.BuildPlayerGameLogRequest("1628389"))
	want := "https://stats.nba.com/stats/playergamelog?PlayerID=1628389&Season=2024-25&SeasonType=Regular+Season"

	if !urlsEqual(t, want, got) {
		t.Errorf("BuildPlayerGameLogRequest() got %s, want %s", got, want)
	}
}

func TestClient_FetchPlayerProfile(t *testing.T) {
	tests := []struct {
		name         string
		httpResponse []byte
		httpErr      error
		writeFileErr error
		fileExists   bool
		wantErr      bool
	}{
		{
			name:         "successful fetch all 3 endpoints",
			httpResponse: []byte(`{"data":"success"}`),
			fileExists:   false,
			wantErr:      false,
		},
		{
			name:       "cached files skip fetch",
			fileExists: true,
			wantErr:    false,
		},
		{
			name:         "http error logged but not returned",
			httpResponse: nil,
			httpErr:      errors.New("http error"),
			fileExists:   true, // skip file-write count check
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fetchedPaths []string

			mockHTTP := &MockHTTPClient{
				getFunc: func(url RequestURL) ([]byte, error) {
					return tt.httpResponse, tt.httpErr
				},
			}

			mockRequestBuilder := &MockRequestBuilder{}

			mockPaths := &MockPathManager{
				fullPathFunc: func(name, param string) string {
					return fmt.Sprintf("/tmp/nba/%s_%s", name, param)
				},
			}

			mockFS := &MockFileSystem{
				fileExistsFunc: func(path string) bool {
					return tt.fileExists
				},
				writeFileFunc: func(path string, data []byte) error {
					fetchedPaths = append(fetchedPaths, path)
					return tt.writeFileErr
				},
			}

			client := &Client{
				http:       mockHTTP,
				requests:   mockRequestBuilder,
				Paths:      mockPaths,
				FileSystem: mockFS,
			}

			err := client.FetchPlayerProfile("1628389")
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchPlayerProfile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !tt.fileExists && len(fetchedPaths) != 3 {
				t.Errorf("expected 3 files written, got %d", len(fetchedPaths))
			}

			if tt.fileExists && len(fetchedPaths) != 0 {
				t.Errorf("expected 0 files written (cached), got %d", len(fetchedPaths))
			}
		})
	}
}
