// Make balanced rosters according to weighted criteria

package main

import (
	"flag"
	"fmt"
)

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
