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
	// form (player name, value, gender ('m' or 'f'))
	playersCsv := csv.NewReader(file)
	_, err = playersCsv.Read()
	baseutil.Check(err)

	// Read in all players
	playersCsvLines, err := playersCsv.ReadAll()
	players := make([]Player, len(playersCsvLines))
	baseutil.Check(err)
	for i, player := range playersCsvLines {
		value, err := strconv.ParseFloat(player[1], 32)
		baseutil.Check(err)
		gender, err := StringToGender(player[2])
		baseutil.Check(err)
		players[i] = Player{player[0], float32(value), gender, 0}
	}
	return players
}
