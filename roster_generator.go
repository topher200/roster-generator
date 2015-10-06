// Make balanced rosters according to weighted criteria

package main

import (
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/GaryBoone/GoStats/stats"
	"github.com/topher200/baseutil"
)

const numTeams = 6

type Solution struct {
	players []Player
}

type Team struct {
	players []Player
}

func splitIntoTeams(players []Player) []Team {
	teams := make([]Team, numTeams)
	for _, player := range players {
		teams[player.team].players = append(teams[player.team].players, player)
	}
	return teams
}

// Score a solution based on weighted critera.
func score(solution Solution) float64 {
	teams := splitIntoTeams(solution.players)

	// Balanced by number
	teamLengths := make([]float64, numTeams)
	for i, team := range teams {
		teamLengths[i] = float64(len(team.players))
	}
	fmt.Println("teamLengths", teamLengths)
	teamsStdDev := stats.StatsSampleStandardDeviation(teamLengths)
	fmt.Println("teamsStdDev", teamsStdDev)

	// Score on balance in gender.
	//
	// For each Gender we make a list of the number of players of that gender on
	// each team. Then we take the standard deviation of those two lists to
	// determine the gender imbalance.
	teamGenders := make(map[Gender][]int)
	for _, gender := range []Gender{Male, Female} {
		teamGenders[gender] = make([]int, 6)
	}
	for teamNum, team := range teams {
		for _, player := range team.players {
			teamGenders[player.gender][teamNum] += 1
		}
	}
	fmt.Println("teamGenders", teamGenders)
	for gender, teamList := range teamGenders {
		teamsStdDev = baseutil.StandardDeviationInt(teamList)
		fmt.Println("gender", gender, "std dev:", teamsStdDev)
	}

	return teamsStdDev
}

func main() {
	// Read command line input
	filenamePointer := kingpin.Arg("input-file",
		"filename from which to get list of players").
		Required().String()
	deterministicPointer := kingpin.Flag("deterministic",
		"makes our output deterministic by allowing the default rand.Seed").
		Short('d').Bool()
	kingpin.Parse()

	if !*deterministicPointer {
		rand.Seed(time.Now().UTC().UnixNano())
	}

	players := ParsePlayers(*filenamePointer)
	solution := Solution{players}

	fmt.Println("score:", score(solution))
}
