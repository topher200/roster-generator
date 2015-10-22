// Make balanced rosters according to weighted criteria

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/pkg/profile"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Genetic algorithm constants
const (
	// Number of teams to break players into
	numTeams = 6
	// Number of times to run our genetic algorithm
	numRuns = 100
	// Percent of the time we will try to mutate. After each
	// mutation, we have a mutationChance percent chance of
	// mutating again.
	mutationChance = 5
	// We will make numSolutionsPerRun every run, and numParents carry
	// over into the next run to create the next batch of solutions.
	numSolutionsPerRun = 1000
	numParents         = 20
)

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
		fmt.Printf("Team #%d, %d players. Average rating: %.2f\n",
			i, len(teams[i].players), AverageRating(team))
		writer := new(tabwriter.Writer)
		writer.Init(os.Stdout, 0, 0, 1, ' ', 0)
		for _, filterFunc := range []PlayerFilter{IsMale, IsFemale} {
			filteredPlayers := Filter(team.players, filterFunc)
			sort.Sort(sort.Reverse(ByRating(filteredPlayers)))
			for _, player := range filteredPlayers {
				fmt.Fprintln(writer, player)
			}
		}
		writer.Flush()
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

	solutionScore, _ := ScoreSolution(newPlayers)
	return Solution{newPlayers, solutionScore}
}

// performRun creates a new solution list by breeding parents.
func performRun(parents []Solution) []Solution {
	newSolutions := make([]Solution, numSolutionsPerRun)
	for i, _ := range newSolutions {
		if i < numParents {
			// Keep the parents from last time - elitism!
			newSolutions[i] = parents[i]
		} else {
			// Make a new solution based on two random parents
			newSolutions[i] = breed(
				parents[rand.Intn(len(parents))],
				parents[rand.Intn(len(parents))])
		}
	}
	return newSolutions
}

// parseCommandLine returns a list of players and a bool for running profiler
func parseCommandLine() ([]Player, bool) {
	filenamePointer := kingpin.Arg("input-file",
		"filename from which to get list of players").
		Required().String()
	deterministicPointer := kingpin.Flag("deterministic",
		"makes our output deterministic by allowing the default rand.Seed").
		Short('d').Bool()
	runProfilingPointer := kingpin.Flag("profiling",
		"output profiling stats when true").Short('p').Bool()
	kingpin.Parse()

	if !*deterministicPointer {
		rand.Seed(time.Now().UTC().UnixNano())
	}

	return ParsePlayers(*filenamePointer), *runProfilingPointer
}

func main() {
	players, profilingOn := parseCommandLine()
	if len(players) == 0 {
		panic("Could not find players")
	}

	// Start profiler
	if profilingOn {
		log.Println("Running profiler")
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}

	// Create random Parent solutions to start
	parentSolutions := make([]Solution, numParents)
	for i, _ := range parentSolutions {
		ourPlayers := make([]Player, len(players))
		copy(ourPlayers, players)
		randomizeTeams(ourPlayers)
		solutionScore, _ := ScoreSolution(ourPlayers)
		parentSolutions[i] = Solution{ourPlayers, solutionScore}
	}

	// Use the random starting solutions to determine the worst case for each of
	// our criteria
	PopulateWorstCases(parentSolutions)

	topScore := parentSolutions[0].score
	for i := 0; i < numRuns; i++ {
		// If we have a new best score, save and print it!
		if topScore != parentSolutions[0].score {
			topScore = parentSolutions[0].score
			log.Println("New top score! Run number ", i, "Score:", topScore)
			PrintSolutionScoring(parentSolutions[0])
		}

		// Create new solutions, and save the best ones
		newSolutions := performRun(parentSolutions)
		sort.Sort(ByScore(newSolutions))
		for i, _ := range parentSolutions {
			parentSolutions[i] = newSolutions[i]
		}
	}

	// Display our solution to the user
	topSolution := parentSolutions[0]
	log.Printf("Top score is %f, solution: %v\n", topSolution, topSolution)
	PrintTeams(topSolution)
	PrintSolutionScoring(topSolution)
}
