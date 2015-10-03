// Make balanced rosters according to weighted criteria

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"unicode/utf8"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Player struct {
	name   string
	value  float32
	gender rune
	team   uint8
}

type Solution struct {
	players []Player
}

func ParsePlayers(inputFilename string) []Player {
	// Open our input Players file
	file, err := os.Open(inputFilename)
	check(err)
	defer file.Close()

	// Read in our csv. Throw away the header. We expect the input to be of the
	// form (player name, value, gender ('m' or 'f'))
	playersCsv := csv.NewReader(file)
	_, err = playersCsv.Read()
	check(err)

	// Read in all players
	playersCsvLines, err := playersCsv.ReadAll()
	players := make([]Player, len(playersCsvLines))
	check(err)
	for i, player := range playersCsvLines {
		value, err := strconv.ParseFloat(player[1], 32)
		check(err)
		gender, _ := utf8.DecodeRuneInString(player[2])
		players[i] = Player{player[0], float32(value), gender, 0}
	}
	return players
}

func main() {
	// Read command line input
	filenamePointer := flag.String(
		"input-file", "input-test.txt", "filename to read input")
	flag.Parse()

	playerList := ParsePlayers(*filenamePointer)

	fmt.Println(playerList)
}
