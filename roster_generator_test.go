package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitIntoTeams(t *testing.T) {
	// todo test remove 2
	players := make([]Player, 2)

	players[0] = Player{Name{"Team 1", "Player"}, 100, Male, 1}
	players[1] = Player{Name{"Team 2", "Player"}, 100, Male, 2}

	teams := splitIntoTeams(players)

	assert.Equal(t, 6, len(teams))
	assert.Equal(t, 0, len(teams[0].players))
	assert.Equal(t, 1, len(teams[1].players))
	assert.Equal(t, 1, len(teams[2].players))
}
