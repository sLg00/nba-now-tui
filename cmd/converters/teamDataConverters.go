package converters

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"log"
)

func PopulateTeamInfo(rs types.ResponseSet) (types.TeamCommonInfo, []string, error) {
	response := rs
	headers := response.ResultSets[0].Headers
	if len(headers) == 0 {
		return types.TeamCommonInfo{}, nil, fmt.Errorf("could not unmarshall team info: no %v returned", headers)
	}

	var currentTeam types.TeamCommonInfo
	for _, row := range response.ResultSets[0].RowSet {
		if len(row) != len(headers) {
			err := fmt.Errorf("header length mismatch with row: %v vs %v", len(headers), len(row))
			log.Println(err)
			return types.TeamCommonInfo{}, nil, err
		}
		for i, value := range row {
			header := headers[i]
			switch header {
			case "TEAM_ID":
				if v, ok := value.(float64); ok {
					currentTeam.TeamID = int(v)
				}
			case "SEASON_YEAR":
				if v, ok := value.(string); ok {
					currentTeam.SeasonYear = v
				}
			case "TEAM_CITY":
				if v, ok := value.(string); ok {
					currentTeam.TeamCity = v
				}
			case "TEAM_NAME":
				if v, ok := value.(string); ok {
					currentTeam.TeamName = v
				}
			case "TEAM_ABBREVIATION":
				if v, ok := value.(string); ok {
					currentTeam.TeamAbbrev = v
				}
			case "TEAM_CONFERENCE":
				if v, ok := value.(string); ok {
					currentTeam.TeamConference = v
				}
			case "TEAM_DIVISION":
				if v, ok := value.(string); ok {
					currentTeam.TeamDivision = v
				}
			case "TEAM_CODE":
				if v, ok := value.(string); ok {
					currentTeam.TeamCode = v
				}
			case "TEAM_SLUG":
				if v, ok := value.(string); ok {
					currentTeam.TeamSlug = v
				}
			case "W":
				if v, ok := value.(float64); ok {
					currentTeam.Wins = int(v)
				}
			case "L":
				if v, ok := value.(float64); ok {
					currentTeam.Losses = int(v)
				}
			case "PCT":
				if v, ok := value.(float64); ok {
					currentTeam.WinPct = v
				}
			case "CONF_RANK":
				if v, ok := value.(float64); ok {
					currentTeam.ConfRank = int(v)
				}
			case "DIV_RANK":
				if v, ok := value.(float64); ok {
					currentTeam.DivRank = int(v)
				}
			case "MIN_YEAR":
				if v, ok := value.(float64); ok {
					currentTeam.MinYear = int(v)
				}
			case "MAX_YEAR":
				if v, ok := value.(float64); ok {
					currentTeam.MaxYear = int(v)
				}
			}
		}
	}
	return currentTeam, headers, nil
}

// PopulateTeamStats maps the data to Teams struct to display season standings
func PopulateTeamStats(rs types.ResponseSet) (types.Teams, []string, error) {
	response := rs
	headers := response.ResultSets[0].Headers
	if len(headers) == 0 {
		return nil, nil, fmt.Errorf("could not unmarshall team stats: no %v returned", headers)
	}

	var teamStats types.Teams
	for _, row := range response.ResultSets[0].RowSet {
		if len(row) != len(headers) {
			err := fmt.Errorf("header row length does not match row length: %v != %v", len(headers), len(row))
			log.Println(err)
			return nil, nil, err
		}
		var team types.Team
		for i, value := range row {
			switch v := value.(type) {
			case float64:
				switch headers[i] {
				case headers[2]:
					team.TeamID = int(v)
				case headers[8]:
					team.PlayoffRank = int(v)
				case headers[12]:
					team.DivisionRank = int(v)
				case headers[13]:
					team.Wins = int(v)
				case headers[14]:
					team.Losses = int(v)
				case headers[16]:
					team.LeagueRank = int(v)
				case headers[24]:
					team.LongHomeStreak = int(v)
				case headers[26]:
					team.LongRoadStreak = int(v)
				case headers[28]:
					team.LongWinStreak = int(v)
				case headers[30]:
					team.LongLossStreak = int(v)
				case headers[32]:
					team.CurrentHomeStreak = int(v)
				case headers[34]:
					team.CurrentRoadStreak = int(v)
				case headers[36]:
					team.CurrentStreak = int(v)
				case headers[39]:
					team.ClinchedConferenceTitle = int(v)
				case headers[40]:
					team.ClinchedDivisionTitle = int(v)
				case headers[41]:
					team.ClinchedPlayoffBirth = int(v)
				case headers[42]:
					team.ClinchedPlayIn = int(v)
				case headers[43]:
					team.EliminatedConference = int(v)
				case headers[44]:
					team.EliminatedDivision = int(v)
				case headers[74]:
					team.TotalPoints = int(v)
				case headers[75]:
					team.OppTotalPoints = int(v)
				case headers[76]:
					team.DiffTotalPoints = int(v)
				case headers[78]:
					team.PlayoffSeeding = int(v)
				case headers[79]:
					team.ClinchedPostSeason = int(v)
				case headers[15]:
					team.WinPCT = v
				case headers[63]:
					team.PointsPG = v
				case headers[64]:
					team.OppPointsPG = v
				case headers[65]:
					team.DiffPointsPG = v
				case headers[37]:
					team.ConferenceGamesBack = v
				case headers[38]:
					team.DivisionGamesBack = v
				case headers[77]:
					team.LeagueGamesBack = v
				}
			case string:
				switch headers[i] {
				case headers[0]:
					team.LeagueID = v
				case headers[1]:
					team.SeasonID = v
				case headers[3]:
					team.TeamCity = v
				case headers[4]:
					team.TeamName = v
				case headers[5]:
					team.TeamSlug = v
				case headers[6]:
					team.Conference = v
				case headers[7]:
					team.ConferenceRecord = v
				case headers[9]:
					team.ClinchIndicator = v
				case headers[10]:
					team.Division = v
				case headers[11]:
					team.DivisionRecord = v
				case headers[17]:
					team.Record = v
				case headers[18]:
					team.Home = v
				case headers[19]:
					team.Road = v
				case headers[20]:
					team.L10 = v
				case headers[21]:
					team.Last10Home = v
				case headers[22]:
					team.Last10Road = v
				case headers[23]:
					team.OT = v
				case headers[25]:
					team.ThreePTSOrLess = v
				case headers[27]:
					team.TenPTSOrMore = v
				case headers[29]:
					team.StrLongHomeStreak = v
				case headers[31]:
					team.StrLongRoadStreak = v
				case headers[33]:
					team.StrCurrentHomeStreak = v
				case headers[35]:
					team.StrCurrentRoadStreak = v
				case headers[59]:
					team.VsEast = v
				case headers[60]:
					team.VsAtlantic = v
				case headers[61]:
					team.VsCentral = v
				case headers[62]:
					team.VsSoutheast = v
				case headers[66]:
					team.VsWest = v
				case headers[67]:
					team.VsNorthwest = v
				case headers[68]:
					team.VsPacific = v
				case headers[69]:
					team.VsSouthwest = v
				case headers[70]:
					team.Jan = v
				case headers[71]:
					team.Feb = v
				case headers[72]:
					team.Mar = v
				case headers[73]:
					team.Apr = v
				case headers[74]:
					team.May = v
				case headers[75]:
					team.Jun = v
				case headers[76]:
					team.Jul = v
				case headers[77]:
					team.Aug = v
				case headers[78]:
					team.Sep = v
				case headers[79]:
					team.Oct = v
				case headers[80]:
					team.Nov = v
				case headers[81]:
					team.Dec = v
				case headers[82]:
					team.Score80Plus = v
				case headers[83]:
					team.OppScore80Plus = v
				case headers[84]:
					team.ScoreBelow80 = v
				case headers[85]:
					team.OppScoreBelow80 = v
				}

			}
		}
		teamStats = append(teamStats, team)
	}
	return teamStats, headers, nil
}
