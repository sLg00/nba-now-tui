package logos

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed *
var logoFiles embed.FS

// LoadTeamLogo reads the contents of the logo directory and returns a logo based on the input parameter value (team name)
func LoadTeamLogo(s string) string {
	logos := make(map[string]string)
	logo := "logo here"

	files, err := fs.ReadDir(logoFiles, ".")
	if err != nil {
		log.Println("could not read logos directory")
		return logo
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		teamName := file.Name()
		logoContent, err := logoFiles.ReadFile(teamName)

		if err != nil {
			log.Printf("could not read logo %s: %v", teamName, err)
			continue
		}

		logos[teamName] = string(logoContent)
		logo = logos[s]
	}

	return logo
}
