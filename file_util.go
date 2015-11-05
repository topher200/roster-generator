// File operations. Retrieve players from csv

package main

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/topher200/baseutil"
)

func ParsePlayers(inputFilename string) []Player {
	// Open our input Players file
	file, err := os.Open(inputFilename)
	baseutil.Check(err)
	defer file.Close()

	// Read in our csv. Throw away the header. Because we're getting our input
	// directly from the league signup form, we expect the input to be shaped like
	// this:
	// Field 3: First name
	// Field 4: Last name
	// Field 8: "Male" or "Female"
	// Field 38: Rating
	playersCsv := csv.NewReader(file)
	columnNames, err := playersCsv.Read()
	baseutil.Check(err)

	// Read in all players
	playersCsvLines, err := playersCsv.ReadAll()
	baseutil.Check(err)
	players := make([]Player, len(playersCsvLines))
	rows := make([]map[string]string, len(playersCsvLines))
	for rowNum, row := range playersCsvLines {
		rows[rowNum] = make(map[string]string)
		for columnNum, value := range row {
			if value == "" {
				continue
			}
			rows[rowNum][columnNames[columnNum]] = value
		}
	}

	for i, row := range rows {
		firstName := row["firstname"]
		lastName := row["lastname"]
		gender, err := StringToGender(row["gender"])
		baseutil.Check(err)
		rating, err := strconv.ParseFloat(row["rating"], 32)
		baseutil.Check(err)
		players[i] = Player{
			Name{firstName, lastName}, float32(rating), gender, uint8(0), Name{}}
	}
	return players
}

// ParseBaggages has the side effect of setting the .baggage for all Players
func ParseBaggages(inputFilename string, players []Player) {
	// Read in our csv. Throw away the header. We expect this format:
	// Field 0: Player 1 First Name
	// Field 1: Player 1 Last Name
	// Field 2: Player 2 First Name
	// Field 3: Player 2 Last Name
	file, err := os.Open(inputFilename)
	baseutil.Check(err)
	defer file.Close()
	baggagesCsv := csv.NewReader(file)
	_, err = baggagesCsv.Read()
	baseutil.Check(err)

	baggagesCsvLines, err := baggagesCsv.ReadAll()
	baseutil.Check(err)
	for _, baggage := range baggagesCsvLines {
		playerPointer, err := FindPlayer(players, Name{baggage[0], baggage[1]})
		baseutil.Check(err)
		if playerPointer.HasBaggage() {
			newLog.Panicf("Player %v already has baggage %v",
				*playerPointer, playerPointer.baggage)
		}
		playerPointer.baggage = Name{baggage[2], baggage[3]}
		newLog.Info("Found baggage of %v for %v",
			playerPointer.baggage, playerPointer.String())
	}
}
