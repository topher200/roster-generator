// Make balanced rosters according to weighted criteria

package main

import (
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/topher200/baseutil"
)

const numTeams = 6
const numRuns = 10

type Solution struct {
	players []Player
	score   float64
}

// Implement sort.Interface for []Solution, sorting based on score
type ByScore []Solution

func (a ByScore) Len() int {
	return len(a)
}
func (a ByScore) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByScore) Less(i, j int) bool {
	return a[i].score < a[j].score
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

// Weights to use to determine criteria importance
const numberBalance = 10
const genderBalance = 8

// Score a solution based on weighted critera.
func score(players []Player) float64 {
	teams := splitIntoTeams(players)

	// Balanced by number
	teamLengths := make([]int, numTeams)
	for i, team := range teams {
		teamLengths[i] = len(team.players)
	}
	teamsStdDev := baseutil.StandardDeviationInt(teamLengths)

	totalScore := teamsStdDev

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
	for _, teamList := range teamGenders {
		teamsStdDev = baseutil.StandardDeviationInt(teamList)
		totalScore += teamsStdDev
	}

	return totalScore
}

func randomizeTeams(players []Player) {
	for i, _ := range players {
		players[i].team = uint8(rand.Intn(numTeams))
	}
}

// Breed via combining the two given solutions, then randomly mutating.
func breed(solution1 Solution, solution2 Solution) Solution {
	// TODO
	return solution1
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
	if len(players) == 0 {
		panic("Could not find players")
	}

	// Create two random solutions to start
	topSolutions := make([]Solution, 2)
	for i, _ := range topSolutions {
		randomizeTeams(players)
		topSolutions[i] = Solution{players, score(players)}
	}

	for i := 0; i < numRuns; i++ {
		fmt.Println("top scores:", topSolutions[0].score, topSolutions[1].score)

		// Create new solutions by breeding the top two solutions
		newSolutions := make([]Solution, 20)
		for i, _ := range newSolutions {
			// Keep the top solutions from last time - elitism!
			if i <= 1 {
				newSolutions[i] = topSolutions[i]
				continue
			}
			newSolutions[i] = breed(topSolutions[0], topSolutions[1])
		}

		// Of all the solutions we now have, save only our best two
		sortedSolutions := ByScore(newSolutions)
		topSolutions[0], topSolutions[1] = sortedSolutions[0], sortedSolutions[1]
	}

	fmt.Println("top scores:", topSolutions[0].score, topSolutions[1].score)
}
