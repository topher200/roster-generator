// Data structs (and functions to act on them) to hold Player information

package main

import "fmt"

type Gender uint8

const (
	Male Gender = iota
	Female
	Default
)

// StringToGender parses single character string into a Gender.
//
// The string must be either "m" or "f", or we return error.
func StringToGender(s string) (Gender, error) {
	switch s {
	case "m":
		return Male, nil
	case "f":
		return Female, nil
	}
	return Default, fmt.Errorf("invalid gender '%s'", s)
}

func IsMale(player Player) bool {
	return player.gender == Male
}
func IsFemale(player Player) bool {
	return player.gender == Female
}

type Player struct {
	name   string
	rating float32
	gender Gender
	team   uint8
}

// Implement fmt.Stringer for printing players
func (player Player) String() string {
	return fmt.Sprintf(" %v: rating %v\n", player.name, player.rating)
}
