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

	// Read in our csv. Throw away the header. We expect the input to be of the
	// form:
	// Field 1: First name
	// Field 2: Last name
	// Field 6: "Male" or "Female"
	// Field 33: Rating
	playersCsv := csv.NewReader(file)
	_, err = playersCsv.Read()
	baseutil.Check(err)

	// Read in all players
	playersCsvLines, err := playersCsv.ReadAll()
	players := make([]Player, len(playersCsvLines))
	baseutil.Check(err)
	for i, player := range playersCsvLines {
		firstName := player[1]
		lastName := player[2]
		gender, err := StringToGender(player[6])
		baseutil.Check(err)
		rating, err := strconv.ParseFloat(player[33], 32)
		baseutil.Check(err)
		players[i] = Player{
			firstName, lastName, float32(rating), gender, uint8(0)}
	}
	return players
}
