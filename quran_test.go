package quranize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAyaFound(t *testing.T) {
	text, err := QuranSimpleClean.GetAya(1, 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, text)
}

func TestGetAyaSuraNotFound(t *testing.T) {
	_, err := QuranSimpleClean.GetAya(0, 0)
	assert.Error(t, err)
}

func TestGetAyaAyaNotFound(t *testing.T) {
	_, err := QuranSimpleClean.GetAya(1, 0)
	assert.Error(t, err)
}

func TestGetSuraNameFound(t *testing.T) {
	text, err := QuranSimpleClean.GetSuraName(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, text)
}

func TestGetSuraNameNotFound(t *testing.T) {
	_, err := QuranSimpleClean.GetSuraName(0)
	assert.Error(t, err)
}

func TestLocateEmptyString(t *testing.T) {
	input := ""
	expected := []Location{}
	actual := QuranSimpleClean.Locate(input)
	assert.Equal(t, expected, actual)
}

func TestLocateNonAlquran(t *testing.T) {
	input := "alfan"
	expected := []Location{}
	actual := QuranSimpleClean.Locate(input)
	assert.Equal(t, expected, actual)
}

func TestLocateAlquran(t *testing.T) {
	input := "بسم الله الرحمن الرحيم"
	expected := []Location{Location{1, 1, 0}, Location{27, 30, 4}}
	actual := QuranSimpleClean.Locate(input)
	assert.Equal(t, expected, actual)
}
