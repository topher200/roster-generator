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
	firstName string
	lastName  string
	rating    float32
	gender    Gender
	team      uint8
}

// FindPlayer returns the first matching player in the list of players.
//
// Return error if none are found
func FindPlayer(players []Player, firstName string, lastName string) (
	Player, error) {
	for _, player := range players {
		if player.firstName == firstName && player.lastName == lastName {
			return player, nil
		}
	}
	return Player{}, fmt.Errorf(
		"No player with name '%s' '%s' found", firstName, lastName)
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
	return fmt.Sprintf("%v %v,\t%v,\trating: %v",
		player.firstName, player.lastName, genderString, player.rating)
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
