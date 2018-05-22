package quran

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadTransliterationFileNotFound(t *testing.T) {
	assert.Panics(t, func() { loadTransliteration("") })
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
