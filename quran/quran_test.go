package quran

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAyaFound(t *testing.T) {
	text, err := QuranClean.GetAya(1, 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, text)
}

func TestGetAyaSuraNotFound(t *testing.T) {
	_, err := QuranClean.GetAya(0, 0)
	assert.Error(t, err)
}

func TestGetAyaAyaNotFound(t *testing.T) {
	_, err := QuranClean.GetAya(1, 0)
	assert.Error(t, err)
}

func TestGetSuraNameFound(t *testing.T) {
	text, err := QuranClean.GetSuraName(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, text)
}

func TestGetSuraNameNotFound(t *testing.T) {
	_, err := QuranClean.GetSuraName(0)
	assert.Error(t, err)
}

func TestLocateEmptyString(t *testing.T) {
	input := ""
	expected := []Location{}
	actual := QuranClean.Locate(input)
	assert.Equal(t, expected, actual)
}

func TestLocateNonAlquran(t *testing.T) {
	input := "alfan"
	expected := []Location{}
	actual := QuranClean.Locate(input)
	assert.Equal(t, expected, actual)
}

func TestLocateAlquran(t *testing.T) {
	input := "بسم الله الرحمن الرحيم"
	expected := []Location{Location{0, 0, 0}, Location{26, 29, 4}}
	actual := QuranClean.Locate(input)
	assert.Equal(t, expected, actual)
}
