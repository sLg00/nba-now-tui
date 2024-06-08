package datamodels

import (
	"github.com/sLg00/nba-now-tui/app/internal/client"
	"log"
)

type Team struct {
	LeagueID                string  `json:"LeagueID"`
	SeasonID                string  `json:"SeasonID"`
	TeamID                  int     `json:"TeamID"`
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
	WinPCT                  float64 `json:"WinPCT"`
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

func PopulateTeamStats() (Teams, []string, error) {
	response, err := unmarshallResponseJSON(client.SSFULLPATH)
	if err != nil {
		log.Println("Could not unmarshall json data:", err)
	}
	headers := response.ResultSets[0].Headers
	var teamStats Teams
	for _, row := range response.ResultSets[0].RowSet {
		if len(row) != len(headers) {
			log.Println("ERR: Header and row length do not match.")
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
