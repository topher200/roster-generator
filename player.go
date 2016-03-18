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

type Name struct {
	firstName, lastName string
}
type Player struct {
	name     Name
	rating   float32
	gender   Gender
	team     uint8
	baggages []Name
}

// FindPlayer returns the first matching player in the list of players.
//
// Return error if none are found
func FindPlayer(players []Player, name Name) (
	*Player, error) {
	for i, player := range players {
		if player.name == name {
			return &players[i], nil
		}
	}
	return &Player{}, fmt.Errorf(
		"No player with name '%s' found", name)
}

// Implement fmt.Stringer for printing players
func (player Player) String() string {
	return fmt.Sprintf("%.02f %s %s",
		player.rating, player.name.firstName, player.name.lastName)
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

type PlayerFilter func(player Player) bool

func Filter(players []Player, filter PlayerFilter) (filteredPlayers []Player) {
	for _, player := range players {
		if filter == nil || filter(player) {
			filteredPlayers = append(filteredPlayers, player)
		}
	}
	return
}
