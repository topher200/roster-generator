package main

import "github.com/topher200/baseutil"

type criterionCalculationFunction func(teams []Team) Score
type playerFilter func(player Player) bool
type Criterion struct {
	name      string                       // human readable name
	calculate criterionCalculationFunction // how to calculate the raw score
	filter    playerFilter                 // cull down to players that match
	weight    int                          // how much weight to give this score
}

var Criteria = [...]Criterion{
	Criterion{"number of players", playerCountDifference, nil, 10},
	Criterion{"number of males", playerCountDifference, IsMale, 9},
	Criterion{"number of females", playerCountDifference, IsFemale, 9},
}

func playerCountDifference(teams []Team) Score {
	teamLengths := make([]int, numTeams)
	for i, team := range teams {
		teamLengths[i] = len(team.players)
	}
	return Score(baseutil.StandardDeviationInt(teamLengths))
}

// runCriterion by filtering the input teams and running the criterion function
func runCriterion(criterion Criterion, teams []Team) (
	rawScore Score, weightedScore Score) {
	filteredTeams := make([]Team, len(teams))
	for i, team := range teams {
		for _, player := range team.players {
			if criterion.filter == nil || criterion.filter(player) {
				filteredTeams[i].players = append(filteredTeams[i].players, player)
			}
		}
	}

	rawScore = criterion.calculate(filteredTeams)
	weightedScore = Score(float64(rawScore) * float64(criterion.weight))
	return rawScore, weightedScore
}

// Score a solution based on weighted criteria.
func ScorePossibleSolution(players []Player) (totalScore Score) {
	teams := splitIntoTeams(players)
	for _, criterion := range Criteria {
		_, weightedScore := runCriterion(criterion, teams)
		totalScore += weightedScore
	}
	return totalScore
}
