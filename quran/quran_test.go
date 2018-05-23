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
