// Data structs (and functions to act on them) to hold Player information

package main

import "fmt"

type Gender uint8

const (
	Male Gender = iota
	Female
	Default
)

// StringToGender parses a raw input string into a Gender.
//
// The string must be either "Male" or "Female", or we return error.
func StringToGender(s string) (Gender, error) {
	switch s {
	case "Male":
		return Male, nil
	case "Female":
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
	genderString := ""
	switch player.gender {
	case Male:
		genderString = "M"
	case Female:
		genderString = "F"
	}
	return fmt.Sprintf("%v,\t%v,\trating: %v",
		player.name, genderString, player.rating)
}

// Implement sorting for []Player based on rating
type ByRating []Player

func (a ByRating) Len() int {
	return len(a)
}
func (a ByRating) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByRating) Less(i, j int) bool {
	return a[i].rating < a[j].rating
}
