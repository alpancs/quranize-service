package quran

import (
	"os"
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

func TestLoadTransliterationFileNotFound(t *testing.T) {
	assert.Panics(t, func() { loadTransliteration("", nil) })
}

func TestLoadQuranFileNotFound(t *testing.T) {
	assert.Panics(t, func() { loadQuran("", nil) })
}

func TestLoadQuranBadXMLFormat(t *testing.T) {
	assert.Panics(t, func() { loadQuran("arabic-to-alphabet", nil) })
}

func TestGetDefaultCorpusPath(t *testing.T) {
	defer os.Setenv("CORPUS_PATH", os.Getenv("CORPUS_PATH"))
	os.Setenv("CORPUS_PATH", "")
	assert.Equal(t, "corpus/", getCorpusPath(), `default corpus path should be "corpus/"`)
}
