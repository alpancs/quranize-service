package quran

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocateEmptyString(t *testing.T) {
	input := ""
	expected := []Location{}
	actual := Locate(input)
	assert.Equal(t, expected, actual)
}

func TestLocateNonAlquran(t *testing.T) {
	input := "alfan"
	expected := []Location{}
	actual := Locate(input)
	assert.Equal(t, expected, actual)
}

func TestLocateAlquran(t *testing.T) {
	input := "بسم الله الرحمن الرحيم"
	expected := []Location{Location{0, 0, 0}, Location{26, 29, 4}}
	actual := Locate(input)
	assert.Equal(t, expected, actual)
}
