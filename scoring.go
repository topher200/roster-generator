package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/GaryBoone/GoStats/stats"
)

type criterionCalculationFunction func(teams []Team) Score
type criterion struct {
	name      string                       // human readable name
	calculate criterionCalculationFunction // how to calculate the raw score
	filter    PlayerFilter                 // cull down to players that match
	// numPlayers reduces the amount of players we analyze from each team.
	// Sometimes used to just grab the top players on the team, for example.
	// Ignored if 0.
	numPlayers int
	weight     int // how much weight to give this score
	// worstCase is calculated at runtime to be the absolute worst score we can
	// see this criterion getting, calculated using random sampling
	worstCase Score
}

var criteriaToScore = [...]criterion{
	criterion{"matching baggages", baggagesMatch, nil, 0, 10000, 0},
	criterion{"number of players", playerCountDifference, nil, 0, 15, 0},
	criterion{"number of males", playerCountDifference, IsMale, 0, 12, 0},
	criterion{"number of females", playerCountDifference, IsFemale, 0, 12, 0},
	criterion{"average rating players", ratingDifference, nil, 0, 8, 0},
	criterion{"average rating males", ratingDifference, IsMale, 0, 7, 0},
	criterion{"average rating females", ratingDifference, IsFemale, 0, 7, 0},
	criterion{"average rating top players", ratingDifference, nil, 3, 20, 0},
	criterion{"average rating top males", ratingDifference, IsMale, 3, 19, 0},
	criterion{"average rating top females", ratingDifference, IsFemale, 19, 7, 0},
	criterion{"std dev of team player ratings", ratingStdDev, nil, 0, 6, 0},
	criterion{"std dev of team male ratings", ratingStdDev, IsMale, 0, 5, 0},
	criterion{"std dev of team female ratings", ratingStdDev, IsFemale, 0, 5, 0},
}

func playerCountDifference(teams []Team) Score {
	// Score increases as the different in team length becomes greater than 1
	min := len(teams[0].players)
	max := len(teams[0].players)
	for _, team := range teams {
		if len(team.players) < min {
			min = len(team.players)
		}
		if len(team.players) > max {
			max = len(team.players)
		}
	}
	diff := max - min
	// Teams can have +1 person from each other (due to odd # of people). Beyond
	// that it's unbalanced.
	score := diff - 1
	if score < 0 {
		score = 0
	}
	return Score(score)
}

func ratingDifference(teams []Team) Score {
	teamAverageRatings := make([]float64, numTeams)
	for i, team := range teams {
		teamAverageRatings[i] = float64(AverageRating(team))
	}
	return Score(stats.StatsSampleStandardDeviation(teamAverageRatings))
}

func ratingStdDev(teams []Team) Score {
	teamRatingsStdDev := make([]float64, numTeams)
	for i, team := range teams {
		if len(team.players) < 2 {
			teamRatingsStdDev[i] = 0
			continue
		}
		playerRatings := make([]float64, len(team.players))
		for j, player := range team.players {
			playerRatings[j] = float64(player.rating)
		}
		teamRatingsStdDev[i] = stats.StatsSampleStandardDeviation(playerRatings)
	}
	return Score(stats.StatsSampleStandardDeviation(teamRatingsStdDev))
}

func baggagesMatch(teams []Team) Score {
	score := Score(0)
	for _, team := range teams {
		for _, player := range team.players {
			if !player.HasBaggage() {
				continue
			}
			_, err := FindPlayer(team.players, player.baggage)
			if err != nil {
				// Player desired a baggage, but they're not on the team
				score += 1
			}
		}
	}
	return score
}

func AverageRating(team Team) Score {
	if len(team.players) == 0 {
		return Score(0)
	}
	sum := float32(0.0)
	for _, player := range team.players {
		sum += player.rating
	}
	return Score(sum / float32(len(team.players)))
}

// analyze criterion by filtering the input teams and running the criterion's
// function
func (c criterion) analyze(teams []Team) (
	rawScore Score, normalizedScore Score, weightedScore Score) {
	filteredTeams := make([]Team, len(teams))
	for i, _ := range teams {
		players := Filter(teams[i].players, c.filter)
		// If the max num players to run this criterion on is set and we have at
		// least that many players, filter out all but the top ones
		if c.numPlayers > 0 && len(players) > c.numPlayers {
			sort.Sort(sort.Reverse(ByRating(players)))
			players = players[:c.numPlayers]
		}
		filteredTeams[i].players = players
	}

	rawScore = c.calculate(filteredTeams)
	if c.worstCase != 0 {
		normalizedScore = rawScore / c.worstCase
	} else {
		normalizedScore = rawScore
	}
	weightedScore = normalizedScore * Score(c.weight)
	return rawScore, normalizedScore, weightedScore
}

func maxScore(a, b Score) Score {
	if a > b {
		return a
	} else {
		return b
	}
}

// PopulateWorstCases calculates the worst case of each criterion.
//
// The function has the side effect of filling in the worstCase param for each
// criterion in criteriaToScore.
func PopulateWorstCases(solutions []Solution) {
	for _, solution := range solutions {
		_, rawScores := ScoreSolution(solution.players)
		for i, criterion := range criteriaToScore {
			if math.IsNaN(float64(rawScores[i])) {
				continue
			}
			criteriaToScore[i].worstCase = maxScore(
				criterion.worstCase, rawScores[i])
		}
	}
}

// Score a solution based on all known criteria.
//
// Returns the total score for the solution, as well as the raw score found for
// each of the criteriaToScore.
func ScoreSolution(players []Player) (totalScore Score, rawScores []Score) {
	teams := splitIntoTeams(players)
	rawScores = make([]Score, len(criteriaToScore))
	for i, criterion := range criteriaToScore {
		rawScore, _, weightedScore := criterion.analyze(teams)
		rawScores[i] = rawScore
		totalScore += weightedScore
	}
	return totalScore, rawScores
}

func PrintSolutionScoring(solution Solution) {
	teams := splitIntoTeams(solution.players)
	totalScore := Score(0)
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 0, 1, ' ', 0)
	for _, criterion := range criteriaToScore {
		rawScore, normalizedScore, weightedScore := criterion.analyze(teams)
		totalScore += weightedScore
		fmt.Fprintf(
			writer,
			"Balancing %s.\tScore: %.02f\t(= normalized score %.02f * weight %d)\t(raw score %0.2f, worst case %.02f)\tRunning total: %.02f\n",
			criterion.name, weightedScore, normalizedScore, criterion.weight,
			rawScore, criterion.worstCase, totalScore)
	}
	writer.Flush()

	// Print the missing baggages
	for _, team := range teams {
		for _, player := range team.players {
			if !player.HasBaggage() {
				continue
			}
			_, err := FindPlayer(team.players, player.baggage)
			if err != nil {
				// Player desired a baggage, but they're not on the team
				fmt.Printf("%v and %v were unfulfilled baggage\n", player, player.baggage)
			}
		}
	}
}
