// Data structs (and functions to act on them) to hold Player information

package main

import "fmt"

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
