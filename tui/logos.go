package tui

import (
	"log"
	"os"
	"path/filepath"
)

// LoadTeamLogo reads the contents of the logo directory and returns a logo based on the input parameter value (team name)
func LoadTeamLogo(s string) string {
	logos := make(map[string]string)
	logo := "logo here"

	files, err := os.ReadDir("../assets/logos")
	if err != nil {
		log.Println("could not read logos directory")
		return logo
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		teamName := file.Name()
		logoPath := filepath.Join("../assets/logos", teamName)
		logoContent, err := os.ReadFile(logoPath)

		if err != nil {
			log.Printf("could not read logo %s: %v", logoPath, err)
			continue
		}

		logos[teamName] = string(logoContent)
		logo = logos[s]
	}

	return logo
}
