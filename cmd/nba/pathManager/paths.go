package pathManager

import (
	"fmt"
	"github.com/sLg00/nba-now-tui/cmd/nba/types"
	"log"
	"os"
)

type PathManager interface {
	GetFullPath(fileType string, id string) string
	GetBasePaths() []string
}

type PathComps struct {
	Home            string //home directory of the current OS user
	Path            string //path to the config directory of the cmd
	LLFile          string //league leaders file name
	SSFile          string //season standings file name
	DSBFile         string //daily scoreboard file name
	BoxScorePath    string //folder to store box scores
	BoxScoreFile    string //box score file name
	BoxScoreID      string //id of specific box score
	TeamProfilePath string //folder to store profile pages
	//TeamProfileFile string // team profile file
	TeamProfileID   string //id of specific team
	TeamPlayersPath string //folder to store player records
	//TeamPlayersFile string //players of a team file
	NewsCachePath string
	NewsCacheFile string
}

func PathFactory(dates types.DateProvider, id string) PathManager {

	home, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("could not determine home directory: %w", err)
		log.Println(err)
	}

	today, err := dates.GetCurrentDate()
	if err != nil {
		err = fmt.Errorf("could not determine current date: %w", err)
	}

	return &PathComps{
		Home:            home,
		Path:            "/.config/nba-tui/",
		LLFile:          today + "_ll",
		SSFile:          today + "_ss",
		DSBFile:         today + "_dsb",
		BoxScorePath:    "boxscores/",
		BoxScoreFile:    today + "_",
		BoxScoreID:      id,
		TeamProfilePath: "teamprofiles/",
		TeamProfileID:   id,
		TeamPlayersPath: "teamplayers/",
		NewsCachePath:   "news/",
		NewsCacheFile:   today + "_news",
	}
}

func PathFactoryForDate(date string) PathManager {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(fmt.Errorf("could not determine home directory: %w", err))
	}

	return &PathComps{
		Home:            home,
		Path:            "/.config/nba-tui/",
		LLFile:          date + "_ll",
		SSFile:          date + "_ss",
		DSBFile:         date + "_dsb",
		BoxScorePath:    "boxscores/",
		BoxScoreFile:    date + "_",
		TeamProfilePath: "teamprofiles/",
		TeamPlayersPath: "teamplayers/",
		NewsCachePath:   "news/",
		NewsCacheFile:   date + "_news",
	}
}

func (p *PathComps) GetFullPath(fileType string, id string) string {
	base := p.Home + p.Path
	switch fileType {
	case "leagueLeaders":
		return base + p.LLFile
	case "seasonStandings":
		return base + p.SSFile
	case "dailyScores":
		return base + p.DSBFile
	case "boxScore":
		return base + p.BoxScorePath + p.BoxScoreFile + id
	case "teamInfo":
		return base + p.TeamProfilePath + id
	case "playerIndex":
		return base + p.TeamPlayersPath + id
	case "newsCachePath":
		return base + p.NewsCachePath
	case "newsCacheFile":
		return base + p.NewsCachePath + p.NewsCacheFile
	default:
		return base
	}
}

func (p *PathComps) GetBasePaths() []string {
	return []string{
		p.Home + p.Path,
		p.Home + p.Path + p.BoxScorePath,
		p.Home + p.Path + p.TeamPlayersPath,
		p.Home + p.Path + p.NewsCachePath,
	}
}
