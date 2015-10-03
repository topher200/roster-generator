// Make balanced rosters according to weighted criteria

package main

import (
	"flag"
	"fmt"
)

type Gender int

const (
	Male Gender = iota
	Female
	Default
)

func StringToGender(s string) (Gender, error) {
	switch s {
	case "m":
		return Male, nil
	case "f":
		return Female, nil
	}
	return Default, fmt.Errorf("invalid gender '%s'", s)
}

type Player struct {
	name   string
	value  float32
	gender Gender
	team   uint8
}

type Solution struct {
	players []Player
}

func main() {
	// Read command line input
	filenamePointer := flag.String(
		"input-file", "input-test.txt", "filename to read input")
	flag.Parse()

	playerList := ParsePlayers(*filenamePointer)

	fmt.Println(playerList)
}
