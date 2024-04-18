package tui

import "github.com/charmbracelet/bubbles/table"

type gameData struct {
	gameID     int
	teamName1  string
	teamName2  string
	teamScore1 int
	teamScore2 int
}

type boxScore struct {
	statLine   gameData
	team1Stats table.Model
	team2Stats table.Model
}

type dailyScores struct {
	games         []gameData
	selectedIndex int
	boxScores     []boxScore
}

type seasonStandings []table.Model
