// File operations. Retrieve players from csv

package main

import (
	"strconv"

	"github.com/topher200/baseutil"
)

func ParsePlayers(inputFilename string) []Player {
	mappedRows := baseutil.MapReader(inputFilename)
	players := make([]Player, len(mappedRows))
	for i, row := range mappedRows {
		firstName := row["First Name"]
		lastName := row["Last Name"]
		gender, err := StringToGender(row["Gender"])
		baseutil.Check(err)
		rating, err := strconv.ParseFloat(row["Balanced Rating"], 32)
		baseutil.Check(err)
		players[i] = Player{
			Name{firstName, lastName}, float32(rating), gender, uint8(0), []Name{}}
	}
	return players
}

// ParseBaggages has the side effect of setting the .baggage for all Players
func ParseBaggages(inputFilename string, players []Player) {
	for _, baggage := range baseutil.MapReader(inputFilename) {
		playerPointer, err := FindPlayer(
			players, Name{baggage["firstname1"], baggage["lastname1"]})
		baseutil.Check(err)
		playerPointer.baggages = append(
			playerPointer.baggages, Name{baggage["firstname2"], baggage["lastname2"]})
		newLog.Debug("Found baggage of %v for %v",
			playerPointer.baggages[len(playerPointer.baggages)-1], playerPointer.String())
	}
}
