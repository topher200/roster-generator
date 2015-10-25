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

type Baggage struct {
	player1, player2 Player
}

func ParseBaggages(inputFilename string, players []Player) []Baggage {
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
	baggages := make([]Baggage, len(baggagesCsvLines))
	for i, baggage := range baggagesCsvLines {
		player1, err := FindPlayer(players, baggage[0], baggage[1])
		baseutil.Check(err)
		player2, err := FindPlayer(players, baggage[2], baggage[3])
		baseutil.Check(err)
		baggages[i] = Baggage{player1, player2}
	}
	newLog.Info("baggages: %v", baggages)
	return baggages
}
