// Make balanced rosters according to weighted criteria

package main

import (
	"flag"
	"fmt"
)

type Player struct {
	name   string
	value  float32
	gender rune
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
