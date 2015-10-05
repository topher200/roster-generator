// Make balanced rosters according to weighted criteria

package main

import (
	"flag"
	"fmt"
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
	// Balanced by gender
	return 0
}

func main() {
	// Read command line input
	filenamePointer := flag.String(
		"input-file", "input-test.txt", "filename to read input")
	flag.Parse()

	playerList := ParsePlayers(*filenamePointer)

	fmt.Println(playerList)
}
