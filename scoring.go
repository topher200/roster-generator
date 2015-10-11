package main

import (
	"log"

	"github.com/topher200/baseutil"
)

type criterionCalculationFunction func(teams []Team) Score
type PlayerFilter func(player Player) bool
type criterion struct {
	name      string                       // human readable name
	calculate criterionCalculationFunction // how to calculate the raw score
	filter    PlayerFilter                 // cull down to players that match
	weight    int                          // how much weight to give this score
}

var criteriaToScore = [...]criterion{
	criterion{"number of players", playerCountDifference, nil, 10},
	criterion{"number of males", playerCountDifference, IsMale, 9},
	criterion{"number of females", playerCountDifference, IsFemale, 9},
}

func playerCountDifference(teams []Team) Score {
	teamLengths := make([]int, numTeams)
	for i, team := range teams {
		teamLengths[i] = len(team.players)
	}
	return Score(baseutil.StandardDeviationInt(teamLengths))
}

func AverageRating(team Team) Score {
	sum := float32(0.0)
	for _, player := range team.players {
		sum += player.rating
	}
	return Score(sum / float32(len(team.players)))
}

func Filter(players []Player, filter PlayerFilter) (filteredPlayers []Player) {
	for _, player := range players {
		if filter == nil || filter(player) {
			filteredPlayers = append(filteredPlayers, player)
		}
	}
	return
}

// runCriterion by filtering the input teams and running the criterion function
func runCriterion(
	c criterion, teams []Team) (rawScore Score, weightedScore Score) {
	filteredTeams := make([]Team, len(teams))
	for i, _ := range teams {
		filteredTeams[i].players = Filter(teams[i].players, c.filter)
	}

	rawScore = c.calculate(filteredTeams)
	weightedScore = Score(float64(rawScore) * float64(c.weight))
	return rawScore, weightedScore
}

// Score a solution based on all known criteria.
func ScoreSolution(players []Player) (totalScore Score) {
	teams := splitIntoTeams(players)
	for _, criterion := range criteriaToScore {
		_, weightedScore := runCriterion(criterion, teams)
		totalScore += weightedScore
	}
	return totalScore
}

func PrintSolutionScoring(solution Solution) {
	teams := splitIntoTeams(solution.players)
	totalScore := Score(0)
	for _, criterion := range criteriaToScore {
		rawScore, weightedScore := runCriterion(criterion, teams)
		totalScore += weightedScore
		log.Printf(
			"Balancing %s. Raw score %f, weighted score %f. Running total: %f\n",
			criterion.name, rawScore, weightedScore, totalScore)
	}
}