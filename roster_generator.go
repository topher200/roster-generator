// Make balanced rosters according to weighted criteria

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/GaryBoone/GoStats/stats"
)

const numTeams = 6

type Solution struct {
	players []Player
}

type Team struct {
	players []Player
}

func splitIntoTeams(players []Player) [][]Player {
	teams := make([][]Player, numTeams)
	for _, player := range players {
		teams[player.team] = append(teams[player.team], player)
	}
	return teams
}

// Score a solution based on weighted critera.
func score(solution Solution) float64 {
	// Balanced by number
	teams := splitIntoTeams(solution.players)

	teamLengths := make([]float64, numTeams)
	for i, team := range teams {
		teamLengths[i] = float64(len(team))
	}
	fmt.Println("teamLengths", teamLengths)
	teamsStdDev := stats.StatsSampleStandardDeviation(teamLengths)
	fmt.Println("teamsStdDev", teamsStdDev)

	// TODO Balanced by gender

	return teamsStdDev
}

func main() {
	// Read command line input
	filenamePointer := flag.String(
		"input-file", "input-test.txt", "filename to read input")
	flag.Parse()

	players := ParsePlayers(*filenamePointer)
	solution := Solution{players}

	fmt.Println("score:", score(solution))
}
