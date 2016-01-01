// Make balanced rosters according to weighted criteria

package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/op/go-logging"
	"github.com/pkg/profile"

	"gopkg.in/alecthomas/kingpin.v2"
)

var newLog = logging.MustGetLogger("")

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

func maxNumberOfPlayersPerTeam(teams []Team) int {
	maxPlayers := 0
	for i := 0; i < math.MaxInt16; i++ {
		works := false
		for _, team := range teams {
			if len(team.players) >= maxPlayers {
				works = true
			}
		}
		if !works {
			break
		}
		maxPlayers += 1
	}
	return maxPlayers
}

func PrintTeams(solution Solution) {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 0, ' ', 0)
	for _, filterFunc := range []PlayerFilter{IsMale, IsFemale} {
		// Print the rating for each team
		filteredPlayers := Filter(solution.players, filterFunc)
		sort.Sort(sort.Reverse(ByRating(filteredPlayers)))
		teams := splitIntoTeams(filteredPlayers)
		string := ""
		for _, team := range teams {
			string += fmt.Sprintf("|Average: %.02f\t", AverageRating(team))
		}
		string += "|"
		fmt.Fprintln(writer, string)

		string = ""
		for _, team := range teams {
			topPlayers := team.players
			if len(topPlayers) > 3 {
				topPlayers = team.players[:3]
			}
			string += fmt.Sprintf("|Top Average: %.02f\t", AverageRating(Team{topPlayers}))
		}
		string += "|"
		fmt.Fprintln(writer, string)

		// Print the players for each team
		numLoops := maxNumberOfPlayersPerTeam(teams)
		for i := 0; i < numLoops; i++ {
			string := ""
			for _, team := range teams {
				if len(team.players) > i {
					string += fmt.Sprintf("|%s\t", team.players[i].String())
				} else {
					string += "|\t"
				}
			}
			string += "|"
			fmt.Fprintln(writer, string)
		}
	}
	writer.Flush()
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

type workerTask struct {
	parent1, parent2 Solution
}

func worker(tasks <-chan workerTask, results chan<- Solution) {
	for task := range tasks {
		results <- breed(task.parent1, task.parent2)
	}
}

// performRun creates a new solution list by breeding parents.
func performRun(
	parents []Solution, tasks chan<- workerTask, results <-chan Solution) []Solution {
	solutions := make([]Solution, numSolutionsPerRun)

	// Keep the parents from last time - elitism!
	for i := 0; i < numParents; i++ {
		solutions[i] = parents[i]
	}

	// Start jobs
	for i := numParents; i < numSolutionsPerRun; i++ {
		tasks <- workerTask{
			parents[rand.Intn(len(parents))], parents[rand.Intn(len(parents))]}
	}

	// Retreive the results of our jobs
	for i := numParents; i < numSolutionsPerRun; i++ {
		solutions[i] = <-results
	}
	return solutions
}

// parseCommandLine parses the user input
//
// Returns:
//  - a []Player of the players from the input file
//  - a bool which tells us whether or not we should be profiling
//  - the number of CPUs to use for goroutines, which is manipulated by "-d"
func parseCommandLine() ([]Player, bool, int) {
	filenamePointer := kingpin.Arg("players",
		"filename from which to get list of players").
		Required().String()
	baggagesPointer := kingpin.Arg("baggages",
		"filename from which to get list of baggages").
		Required().String()
	deterministicPointer := kingpin.Flag("deterministic",
		"makes our output deterministic by allowing the default rand.Seed").
		Short('d').Bool()
	runProfilingPointer := kingpin.Flag("profiling",
		"output profiling stats when true").Short('p').Bool()
	verbosePointer := kingpin.Flag("verbose",
		"verbose output").Short('v').Bool()
	kingpin.Parse()

	// Set up logging
	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))
	if *verbosePointer {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}

	// To run deterministically, we use the default seed and only one goroutine
	numWorkers := runtime.NumCPU()
	if !*deterministicPointer {
		rand.Seed(time.Now().UTC().UnixNano())
	} else {
		newLog.Info("Seeded deterministically")
		numWorkers = 1
	}

	players := ParsePlayers(*filenamePointer)
	ParseBaggages(*baggagesPointer, players)
	return players, *runProfilingPointer, numWorkers
}

func timeToClose(numRunsCompleted, topScoreRunNumber int) bool {
	return numRunsCompleted > topScoreRunNumber+1000
}

func main() {
	players, profilingOn, numWorkers := parseCommandLine()
	startTime := time.Now()
	if len(players) == 0 {
		panic("Could not find players")
	}

	// Start profiler
	if profilingOn {
		newLog.Info("Running profiler")
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

	// Start our worker goroutines
	tasks := make(chan workerTask, numSolutionsPerRun)
	results := make(chan Solution, numSolutionsPerRun)
	for i := 0; i < numWorkers; i++ {
		go worker(tasks, results)
	}
	defer close(tasks)

	topScore := parentSolutions[0].score
	numRunsCompleted := 0
	topScoreRunNumber := 0
	for {
		// If we have a new best score, save and print it!
		if topScore != parentSolutions[0].score {
			topScore = parentSolutions[0].score
			topScoreRunNumber = numRunsCompleted
			if newLog.IsEnabledFor(logging.DEBUG) {
				newLog.Info("New top score! Run number %d. Score: %.02f",
					numRunsCompleted, topScore)
				PrintSolutionScoring(parentSolutions[0])
			}
		}

		// Create new solutions, and save the best ones
		newSolutions := performRun(parentSolutions, tasks, results)
		sort.Sort(ByScore(newSolutions))
		for i, _ := range parentSolutions {
			parentSolutions[i] = newSolutions[i]
		}
		numRunsCompleted += 1
		if timeToClose(numRunsCompleted, topScoreRunNumber) {
			break
		}
	}

	// Display our solution to the user
	topSolution := parentSolutions[0]
	PrintTeams(topSolution)
	PrintSolutionScoring(topSolution)
	newLog.Debug("Program runtime: %.02fs", time.Since(startTime).Seconds())
}
