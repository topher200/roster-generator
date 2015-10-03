package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToGender(t *testing.T) {
	gender, err := StringToGender("m")
	assert.Nil(t, err)
	assert.Equal(t, gender, Male)
	gender, err = StringToGender("f")
	assert.Nil(t, err)
	assert.Equal(t, gender, Female)
	gender, err = StringToGender("")
	assert.NotNil(t, err)
	gender, err = StringToGender("a")
	assert.NotNil(t, err)
	gender, err = StringToGender("asdf")
	assert.NotNil(t, err)
}
