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
	Dates           types.DateProvider
	Home            string //home directory of the current OS user
	Path            string //path to the config directory of the cmd
	LLFile          string //league leaders file name
	SSFile          string //season standings file name
	DSBFile         string //daily scoreboard file name
	BoxScorePath    string //folder to store box scores
	BoxScoreFile    string //box score file name
	TeamProfilePath string //folder to store profile pages
	TeamProfileFile string // team profile file
}

func PathFactory(dates types.DateProvider, id string) *PathComps {

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
		BoxScoreFile:    today + "_" + id,
		TeamProfilePath: "teamprofiles/" + id,
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
		return base + p.BoxScorePath + id
	case "teamProfile":
		return base + p.TeamProfilePath + id
	default:
		return base
	}
}

func (p *PathComps) GetBasePaths() []string {
	return []string{
		p.Home + p.Path,
		p.Home + p.Path + p.BoxScorePath,
		p.Home + p.Path + p.TeamProfilePath,
	}
}
