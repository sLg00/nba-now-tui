package converters

import (
	"encoding/json"
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
)

// PopulatePlayerBio extracts the player bio from the commonplayerinfo API response.
// The response uses resultSets (plural) with the bio data in the first result set's single row.
func PopulatePlayerBio(rs types.ResponseSet) (types.PlayerBio, error) {
	if len(rs.ResultSets) == 0 || len(rs.ResultSets[0].RowSet) == 0 {
		return types.PlayerBio{}, fmt.Errorf("no player bio data found")
	}

	headers := rs.ResultSets[0].Headers
	row := rs.ResultSets[0].RowSet[0]

	if len(row) != len(headers) {
		return types.PlayerBio{}, fmt.Errorf("header length mismatch: %d vs %d", len(headers), len(row))
	}

	playerData := make(map[string]interface{})
	for i, value := range row {
		playerData[headers[i]] = value
	}

	jsonData, err := json.Marshal(playerData)
	if err != nil {
		return types.PlayerBio{}, fmt.Errorf("failed to marshal player bio: %v", err)
	}

	var bio types.PlayerBio
	if err = json.Unmarshal(jsonData, &bio); err != nil {
		return types.PlayerBio{}, fmt.Errorf("failed to unmarshal player bio: %v", err)
	}

	return bio, nil
}

// PopulateSeasonStats extracts career season stats from the playerprofilev2 API response.
// It looks for the result set named "SeasonTotalsRegularSeason".
func PopulateSeasonStats(rs types.ResponseSet) ([]types.SeasonStats, []string, error) {
	var targetSet *types.ResultSet
	for i := range rs.ResultSets {
		if rs.ResultSets[i].Name == "SeasonTotalsRegularSeason" {
			targetSet = &rs.ResultSets[i]
			break
		}
	}

	if targetSet == nil {
		return nil, nil, fmt.Errorf("SeasonTotalsRegularSeason result set not found")
	}

	headers := targetSet.Headers
	var stats []types.SeasonStats

	for _, row := range targetSet.RowSet {
		if len(row) != len(headers) {
			return nil, nil, fmt.Errorf("row length mismatch: %d vs %d", len(row), len(headers))
		}

		data := make(map[string]interface{})
		for i, value := range row {
			data[headers[i]] = value
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal season stats: %v", err)
		}

		var season types.SeasonStats
		if err = json.Unmarshal(jsonData, &season); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal season stats: %v", err)
		}

		stats = append(stats, season)
	}

	return stats, headers, nil
}

// PopulateGameLog extracts the most recent games from the playergamelog API response.
// Returns at most 5 entries (the API returns games in reverse chronological order).
func PopulateGameLog(rs types.ResponseSet) ([]types.GameLogEntry, []string, error) {
	if len(rs.ResultSets) == 0 {
		return nil, nil, fmt.Errorf("no game log data found")
	}

	headers := rs.ResultSets[0].Headers
	rows := rs.ResultSets[0].RowSet

	limit := len(rows)
	if limit > 5 {
		limit = 5
	}

	var entries []types.GameLogEntry
	for _, row := range rows[:limit] {
		if len(row) != len(headers) {
			return nil, nil, fmt.Errorf("row length mismatch: %d vs %d", len(row), len(headers))
		}

		data := make(map[string]interface{})
		for i, value := range row {
			data[headers[i]] = value
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal game log: %v", err)
		}

		var entry types.GameLogEntry
		if err = json.Unmarshal(jsonData, &entry); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal game log: %v", err)
		}

		entries = append(entries, entry)
	}

	return entries, headers, nil
}
