// Make balanced rosters according to weighted criteria

package main

import (
	"flag"
	"fmt"
	"math"
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
func score(solution Solution) int {
	// Balanced by number
	teams := splitIntoTeams(solution.players)

	// We sum the amount of difference (in number of players on the team) between
	// every team and every other team.
	//
	// This is currently incredibly inefficient (O(n^2) (but n is 6!))
	playerDifference := 0
	for _, team := range teams {
		for _, otherTeam := range teams {
			playerDifference += int(math.Abs(float64(len(team) - len(otherTeam))))
		}
	}

	// Balanced by gender

	return playerDifference
}

func main() {
	// Read command line input
	filenamePointer := flag.String(
		"input-file", "input-test.txt", "filename to read input")
	flag.Parse()

	players := ParsePlayers(*filenamePointer)
	solution := Solution{players}

	fmt.Println(score(solution))
}
