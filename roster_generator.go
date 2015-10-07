// Make balanced rosters according to weighted criteria

package main

import (
	"log"
	"math/rand"
	"sort"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/topher200/baseutil"
)

// Number of teams to break players into
const numTeams = 6

// Number of times to run our genetic algorithm
const numRuns = 100000

// Percent of the time we will try to mutate. After each mutation, we have a
// mutationChance percent chance of mutating again.
const mutationChance = 5

// Weights to use to determine criteria importance
const numberBalance = 10
const genderBalance = 8

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

// Mutate the solution by moving random players to random teams, sometimes.
func mutate(players []Player) {
	for {
		// We have mutationChance of mutating. Otherwise, we break out of our loop
		if rand.Intn(100) > mutationChance {
			return
		}
		// Mutation! Move a random player to a random new team
		players[rand.Intn(len(players))].team = uint8(rand.Intn(numTeams))
	}
}

// Breed via combining the two given solutions, then randomly mutating.
func breed(solution1 Solution, solution2 Solution) Solution {
	// Create the new solution by taking crossover from both inputs
	newPlayers := make([]Player, len(solution1.players))
	for i, _ := range newPlayers {
		// Randomly take each player from solution1 or solution2
		if rand.Intn(100) < 50 {
			newPlayers[i] = solution1.players[i]
		} else {
			newPlayers[i] = solution2.players[i]
		}
	}

	// Mutate the new player list
	mutate(newPlayers)

	return Solution{newPlayers, score(newPlayers)}
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
		ourPlayers := make([]Player, len(players))
		copy(ourPlayers, players)
		randomizeTeams(ourPlayers)
		topSolutions[i] = Solution{ourPlayers, score(ourPlayers)}
	}

	topScore := topSolutions[0].score
	for i := 0; i < numRuns; i++ {
		if topScore != topSolutions[0].score {
			topScore = topSolutions[0].score
			log.Println("New top score! Run number ", i, "Score:", topScore)
		}

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
		sort.Sort(ByScore(newSolutions))
		topSolutions[0], topSolutions[1] = newSolutions[0], newSolutions[1]
	}
}
