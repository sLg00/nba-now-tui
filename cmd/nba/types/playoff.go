package types

type PlayoffTeam struct {
	TeamID  string
	Name    string
	Tricode string
	Seed    int
	Wins    int
}

// PlayoffSeries represents one matchup in the bracket.
// Status: "pre" (projected, not started), "active", "complete".
type PlayoffSeries struct {
	SeriesID   string
	Round      int    // 1=First Round, 2=Semis, 3=Conf Finals, 4=Finals
	Conference string // "East", "West", "Finals"
	TopTeam    PlayoffTeam
	BottomTeam PlayoffTeam
	Status     string
}

type PlayoffBracket struct {
	Season string
	Series []PlayoffSeries // length 15 max; ordered by cursor index (see design doc)
}

type PlayoffGame struct {
	GameID      string
	Date        string
	GameNumber  int
	HomeTricode string
	AwayTricode string
	HomeScore   int
	AwayScore   int
	Completed   bool
}

func (s PlayoffSeries) IsComplete() bool {
	return s.TopTeam.Wins == 4 || s.BottomTeam.Wins == 4
}

// Leader returns the tricode of the team leading the series, or "" if tied.
func (s PlayoffSeries) Leader() string {
	if s.TopTeam.Wins > s.BottomTeam.Wins {
		return s.TopTeam.Tricode
	}
	if s.BottomTeam.Wins > s.TopTeam.Wins {
		return s.BottomTeam.Tricode
	}
	return ""
}
