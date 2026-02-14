package filesystemops

import (
	"encoding/json"
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/pathManager"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

// DataLoader interface is used to inject the relevant ResponseSets into the converter functions
type DataLoader interface {
	LoadDailyScoreboard() (types.ResponseSet, error)
	LoadLeagueLeaders() (types.ResponseSet, error)
	LoadSeasonStandings() (types.ResponseSet, error)
	LoadBoxScore(gameID string) (types.ResponseSet, error)
	LoadTeamInfo(teamID string) (types.ResponseSet, error)
	LoadPlayerIndex(teamID string) (types.ResponseSet, error)
	LoadPlayerInfo(playerID string) (types.ResponseSet, error)
	LoadPlayerCareerStats(playerID string) (types.ResponseSet, error)
	LoadPlayerGameLog(playerID string) (types.ResponseSet, error)
}

// nbaDataLoader implements the DataLoader interface
type nbaDataLoader struct {
	fs    FileSystemHandler
	paths pathManager.PathManager
}

// NewDataLoader is a factory function that instantiates a DataLoader
func NewDataLoader(fs FileSystemHandler, paths pathManager.PathManager) DataLoader {
	return &nbaDataLoader{
		fs:    fs,
		paths: paths,
	}
}

func (dl *nbaDataLoader) LoadDailyScoreboard() (types.ResponseSet, error) {
	path := dl.paths.GetFullPath("dailyScores", "")
	return dl.loadAndUnmarshall(path)
}

func (dl *nbaDataLoader) LoadLeagueLeaders() (types.ResponseSet, error) {
	path := dl.paths.GetFullPath("leagueLeaders", "")
	return dl.loadAndUnmarshall(path)
}

func (dl *nbaDataLoader) LoadSeasonStandings() (types.ResponseSet, error) {
	path := dl.paths.GetFullPath("seasonStandings", "")
	return dl.loadAndUnmarshall(path)
}

func (dl *nbaDataLoader) LoadBoxScore(gameID string) (types.ResponseSet, error) {
	path := dl.paths.GetFullPath("boxScore", gameID)
	return dl.loadAndUnmarshall(path)
}

func (dl *nbaDataLoader) LoadTeamInfo(teamID string) (types.ResponseSet, error) {
	path := dl.paths.GetFullPath("teamInfo", teamID)
	return dl.loadAndUnmarshall(path)
}

func (dl *nbaDataLoader) LoadPlayerIndex(teamID string) (types.ResponseSet, error) {
	path := dl.paths.GetFullPath("playerIndex", teamID)
	return dl.loadAndUnmarshall(path)
}

func (dl *nbaDataLoader) LoadPlayerInfo(playerID string) (types.ResponseSet, error) {
	path := dl.paths.GetFullPath("playerInfo", playerID)
	return dl.loadAndUnmarshall(path)
}

func (dl *nbaDataLoader) LoadPlayerCareerStats(playerID string) (types.ResponseSet, error) {
	path := dl.paths.GetFullPath("playerCareerStats", playerID)
	return dl.loadAndUnmarshall(path)
}

func (dl *nbaDataLoader) LoadPlayerGameLog(playerID string) (types.ResponseSet, error) {
	path := dl.paths.GetFullPath("playerGameLog", playerID)
	return dl.loadAndUnmarshall(path)
}

// loadAnUnmarshall method loads a file using the ReadFile function and thn unmarshalls it into a types.ResponseSet
func (dl *nbaDataLoader) loadAndUnmarshall(path string) (types.ResponseSet, error) {
	data, err := dl.fs.ReadFile(path)
	if err != nil {
		return types.ResponseSet{}, fmt.Errorf("failed to load file %s: %w", path, err)
	}

	var response types.ResponseSet
	if err = json.Unmarshal(data, &response); err != nil {
		return types.ResponseSet{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return response, nil
}
