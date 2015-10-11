// Make balanced rosters according to weighted criteria

package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Number of teams to break players into
const numTeams = 6

// Number of times to run our genetic algorithm
const numRuns = 100

// Percent of the time we will try to mutate. After each mutation, we have a
// mutationChance percent chance of mutating again.
const mutationChance = 5

type Score float64
type Solution struct {
	players []Player
	score   Score
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

func randomizeTeams(players []Player) {
	for i, _ := range players {
		players[i].team = uint8(rand.Intn(numTeams))
	}
}

func PrintTeams(solution Solution) {
	teams := splitIntoTeams(solution.players)
	for i, team := range teams {
		fmt.Printf("Team #%d, %d players. Average rating: %v\n",
			i, len(teams[i].players), AverageRating(team))
		for _, filterFunc := range []PlayerFilter{IsMale, IsFemale} {
			for _, player := range Filter(team.players, filterFunc) {
				fmt.Printf(" %s", player)
			}
		}
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

	return Solution{newPlayers, ScoreSolution(newPlayers)}
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
		topSolutions[i] = Solution{ourPlayers, ScoreSolution(ourPlayers)}
	}

	topScore := topSolutions[0].score
	for i := 0; i < numRuns; i++ {
		if topScore != topSolutions[0].score {
			topScore = topSolutions[0].score
			log.Println("New top score! Run number ", i, "Score:", topScore)
			PrintSolutionScoring(topSolutions[0])
		}

		// Create new solutions by breeding the top two solutions
		newSolutions := make([]Solution, 20)
		for i, _ := range newSolutions {
			if i <= 1 {
				// Keep the top solutions from last time - elitism!
				newSolutions[i] = topSolutions[i]
			} else {
				// Make a new solution based on our top two from last time
				newSolutions[i] = breed(topSolutions[0], topSolutions[1])
			}
		}

		// Of all the solutions we now have, save only our best two
		sort.Sort(ByScore(newSolutions))
		topSolutions[0], topSolutions[1] = newSolutions[0], newSolutions[1]
	}
	topSolution := topSolutions[0]
	log.Printf("Top score is %f, solution: %v\n", topSolution, topSolution)
	PrintTeams(topSolution)
	PrintSolutionScoring(topSolution)
}
