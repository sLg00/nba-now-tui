package datamodels

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/client"
	"github.com/sLg00/nba-now-tui/cmd/nba/httpAPI"
	"log"
)

type Team struct {
	LeagueID                string  `json:"LeagueID" isVisible:"false"`
	SeasonID                string  `json:"SeasonID" isVisible:"false"`
	TeamID                  int     `json:"TeamID" isVisible:"false"`
	TeamCity                string  `json:"TeamCity"`
	TeamName                string  `json:"TeamName"`
	TeamSlug                string  `json:"TeamSlug"`
	Conference              string  `json:"Conference"`
	ConferenceRecord        string  `json:"ConferenceRecord"`
	PlayoffRank             int     `json:"PlayoffRank"`
	ClinchIndicator         string  `json:"ClinchIndicator"`
	Division                string  `json:"Division"`
	DivisionRecord          string  `json:"DivisionRecord"`
	DivisionRank            int     `json:"DivisionRank"`
	Wins                    int     `json:"WINS"`
	Losses                  int     `json:"LOSSES"`
	WinPCT                  float64 `json:"WinPCT" percentage:"true"`
	LeagueRank              int     `json:"LeagueRank"`
	Record                  string  `json:"Record"`
	Home                    string  `json:"HOME"`
	Road                    string  `json:"ROAD"`
	L10                     string  `json:"L10"`
	Last10Home              string  `json:"Last10Home"`
	Last10Road              string  `json:"Last10Road"`
	OT                      string  `json:"OT"`
	ThreePTSOrLess          string  `json:"ThreePTSOrLess"`
	TenPTSOrMore            string  `json:"TenPTSOrMore"`
	LongHomeStreak          int     `json:"LongHomeStreak"`
	StrLongHomeStreak       string  `json:"strLongHomeStreak"`
	LongRoadStreak          int     `json:"LongRoadStreak"`
	StrLongRoadStreak       string  `json:"strLongRoadStreak"` //
	LongWinStreak           int     `json:"LongWinStreak"`
	LongLossStreak          int     `json:"LongLossStreak"`
	CurrentHomeStreak       int     `json:"CurrentHomeStreak"`
	StrCurrentHomeStreak    string  `json:"strCurrentHomeStreak"` //
	CurrentRoadStreak       int     `json:"CurrentRoadStreak"`
	StrCurrentRoadStreak    string  `json:"strCurrentRoadStreak"`
	CurrentStreak           int     `json:"CurrentStreak"`
	StrCurrentStreak        string  `json:"strCurrentStreak"`
	ConferenceGamesBack     float64 `json:"ConferenceGamesBack"`
	DivisionGamesBack       float64 `json:"DivisionGamesBack"`
	ClinchedConferenceTitle int     `json:"ClinchedConferenceTitle"`
	ClinchedDivisionTitle   int     `json:"ClinchedDivisionTitle"`
	ClinchedPlayoffBirth    int     `json:"ClinchedPlayoffBirth"`
	ClinchedPlayIn          int     `json:"ClinchedPlayIn"`
	EliminatedConference    int     `json:"EliminatedConference"`
	EliminatedDivision      int     `json:"EliminatedDivision"`
	AheadAtHalf             string  `json:"AheadAtHalf"`
	BehindAtHalf            string  `json:"BehindAtHalf"`
	TiedAtHalf              string  `json:"TiedAtHalf"`
	AheadAtThird            string  `json:"AheadAtThird"`
	BehindAtThird           string  `json:"BehindAtThird"`
	TiedAtThird             string  `json:"TiedAtThird"`
	Score100PTS             string  `json:"Score100PTS"`
	OppScore100PTS          string  `json:"OppScore100PTS"`
	OppOver500              string  `json:"OppOver500"`
	LeadInFGPCT             string  `json:"LeadInFGPCT"`
	LeadInReb               string  `json:"LeadInReb"`
	FewerTurnovers          string  `json:"FewerTurnovers"`
	PointsPG                float64 `json:"PointsPG"`
	OppPointsPG             float64 `json:"OppPointsPG"`
	DiffPointsPG            float64 `json:"DiffPointsPG"`
	VsEast                  string  `json:"vsEast"`
	VsAtlantic              string  `json:"vsAtlantic"`
	VsCentral               string  `json:"vsCentral"`
	VsSoutheast             string  `json:"vsSoutheast"`
	VsWest                  string  `json:"vsWest"`
	VsNorthwest             string  `json:"vsNorthwest"`
	VsPacific               string  `json:"vsPacific"`
	VsSouthwest             string  `json:"vsSouthwest"`
	Jan                     string  `json:"Jan"`
	Feb                     string  `json:"Feb"`
	Mar                     string  `json:"Mar"`
	Apr                     string  `json:"Apr"`
	May                     string  `json:"May"`
	Jun                     string  `json:"Jun"`
	Jul                     string  `json:"Jul"`
	Aug                     string  `json:"Aug"`
	Sep                     string  `json:"Sep"`
	Oct                     string  `json:"Oct"`
	Nov                     string  `json:"Nov"`
	Dec                     string  `json:"Dec"`
	Score80Plus             string  `json:"Score_80_Plus"`
	OppScore80Plus          string  `json:"Opp_Score_80_Plus"`
	ScoreBelow80            string  `json:"Score_Below_80"`
	OppScoreBelow80         string  `json:"Opp_Score_Below_80"`
	TotalPoints             int     `json:"TotalPoints"`
	OppTotalPoints          int     `json:"OppTotalPoints"`
	DiffTotalPoints         int     `json:"DiffTotalPoints"`
	LeagueGamesBack         float64 `json:"LeagueGamesBack"`
	PlayoffSeeding          int     `json:"PlayoffSeeding"`
	ClinchedPostSeason      int     `json:"ClinchedPostSeason"`
}

type Teams []Team

type TeamCommonInfo struct {
	TeamID         int     `json:"TEAM_ID"`
	SeasonYear     string  `json:"SEASON_YEAR"`
	TeamCity       string  `json:"TEAM_CITY"`
	TeamName       string  `json:"TEAM_NAME"`
	TeamAbbrev     string  `json:"TEAM_ABBREV"`
	TeamConference string  `json:"TEAM_CONFERENCE"`
	TeamDivision   string  `json:"TEAM_DIVISION"`
	TeamCode       string  `json:"TEAM_CODE"`
	TeamSlug       string  `json:"TEAM_SLUG"`
	Wins           int     `json:"W"`
	Losses         int     `json:"L"`
	WinPct         float64 `json:"PCT"`
	ConfRank       int     `json:"CONF_RANK"`
	DivRank        int     `json:"DIV_RANK"`
	MinYear        int     `json:"MIN_YEAR"`
	MaxYear        int     `json:"MAX_YEAR"`
}

// ToStringSlice is a method on the TeamCommonInfo type that enables the attributes of the type to be converted to strings
func (ti TeamCommonInfo) ToStringSlice() []string {
	return structToStringSlice(ti)
}

// ToStringSlice is a method on the Team type that enables the attributes of the type to be converted to strings
func (t Team) ToStringSlice() []string {
	return structToStringSlice(t)
}

// ToStringSlice is a method on the Teams type that enables the attributes of type to be converted to strings
func (ts Teams) ToStringSlice() []string {
	return structToStringSlice(ts)
}

func PopulateTeamInfo(s string, unmarshall func(string) (ResponseSet, error)) (TeamCommonInfo, []string, error) {
	pc := client.NewClient().InstantiatePaths(s).TeamProfileFullPath()
	response, err := unmarshall(pc)
	if err != nil {
		err = fmt.Errorf("could not unmarshall team info: %v", err)
		log.Println()
	}

	headers := response.ResultSets[0].Headers
	if len(headers) == 0 {
		return TeamCommonInfo{}, nil, fmt.Errorf("could not unmarshall team info: no %v returned", headers)
	}

	var currentTeam TeamCommonInfo
	for _, row := range response.ResultSets[0].RowSet {
		if len(row) != len(headers) {
			err = fmt.Errorf("header length mismatch with row: %v vs %v", len(headers), len(row))
			log.Println(err)
			return TeamCommonInfo{}, nil, err
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
func PopulateTeamStats(unmarshall func(string) (ResponseSet, error)) (Teams, []string, error) {
	pc := httpAPI.NewNewClient().InstantiatePaths("").SSFullPath()
	response, err := unmarshall(pc)
	if err != nil {
		err = fmt.Errorf("could not unmarshall team stats: %v", err)
		log.Println(err)
	}

	headers := response.ResultSets[0].Headers
	if len(headers) == 0 {
		return nil, nil, fmt.Errorf("could not unmarshall team stats: no %v returned", headers)
	}

	var teamStats Teams
	for _, row := range response.ResultSets[0].RowSet {
		if len(row) != len(headers) {
			err = fmt.Errorf("header row length does not match row length: %v != %v", len(headers), len(row))
			log.Println(err)
			return nil, nil, err
		}
		var team Team
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

// SplitStandingsPerConference is a method on the Teams type that separates the Western and Easter conference teams into
// separate instances of the Teams type
func (ts Teams) SplitStandingsPerConference() (Teams, Teams) {
	teams := ts
	var westTeams Teams
	var eastTeams Teams

	for _, team := range teams {
		if team.Conference == "East" {
			eastTeams = append(eastTeams, team)
		} else if team.Conference == "West" {
			westTeams = append(westTeams, team)
		}
	}
	return eastTeams, westTeams
}
