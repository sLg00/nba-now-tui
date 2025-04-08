package types

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
	TeamID         int     `json:"TEAM_ID" isVisible:"false" isID:"true"`
	SeasonYear     string  `json:"SEASON_YEAR" isVisible:"false" isID:"false"`
	TeamCity       string  `json:"TEAM_CITY" isVisible:"false" isID:"false"`
	TeamName       string  `json:"TEAM_NAME" isVisible:"false" isID:"false"`
	TeamAbbrev     string  `json:"TEAM_ABBREV" isVisible:"false" isID:"false"`
	TeamConference string  `json:"TEAM_CONFERENCE" isVisible:"false" isID:"false"`
	TeamDivision   string  `json:"TEAM_DIVISION" isVisible:"false" isID:"false"`
	TeamCode       string  `json:"TEAM_CODE" isVisible:"false" isID:"false"`
	TeamSlug       string  `json:"TEAM_SLUG" isVisible:"false" isID:"false"`
	Wins           int     `json:"W" isVisible:"true" isID:"false" display:"Wins" width:"20"`
	Losses         int     `json:"L" isVisible:"true" isID:"false" display:"Losses"`
	WinPct         float64 `json:"PCT" percentage:"true" isVisible:"true" isID:"false" display:"Win %"`
	ConfRank       int     `json:"CONF_RANK" isVisible:"true" isID:"false" display:"Conf rank"`
	DivRank        int     `json:"DIV_RANK" isVisible:"true" isID:"false" display:"Div rank"`
	MinYear        string  `json:"MIN_YEAR" isVisible:"false" isID:"false"`
	MaxYear        string  `json:"MAX_YEAR" isVisible:"false" isID:"false"`
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
