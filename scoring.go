package main

import "github.com/topher200/baseutil"

type criterionCalculationFunction func(teams []Team) Score
type playerFilter func(player Player) bool
type criterion struct {
	name      string                       // human readable name
	calculate criterionCalculationFunction // how to calculate the raw score
	filter    playerFilter                 // cull down to players that match
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

// runCriterion by filtering the input teams and running the criterion function
func runCriterion(c criterion, teams []Team) (
	rawScore Score, weightedScore Score) {
	filteredTeams := make([]Team, len(teams))
	for i, team := range teams {
		for _, player := range team.players {
			if c.filter == nil || c.filter(player) {
				filteredTeams[i].players = append(filteredTeams[i].players, player)
			}
		}
	}

	rawScore = c.calculate(filteredTeams)
	weightedScore = Score(float64(rawScore) * float64(c.weight))
	return rawScore, weightedScore
}

// Score a solution based on all known criteria.
func ScorePossibleSolution(players []Player) (totalScore Score) {
	teams := splitIntoTeams(players)
	for _, criterion := range criteriaToScore {
		_, weightedScore := runCriterion(criterion, teams)
		totalScore += weightedScore
	}
	return totalScore
}
