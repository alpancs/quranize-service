package core

import (
	"os"
	"testing"
)

func TestLoadTransliterationFileNotFound(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("loadTransliteration should panic")
		}
	}()
	loadTransliteration("")
}

func TestLoadQuranFileNotFound(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("loadQuran should panic")
		}
	}()
	loadQuran("", nil)
}

func TestLoadQuranBadXMLFormat(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("loadQuran should panic")
		}
	}()
	loadQuran("arabic-to-alphabet", nil)
}

func TestGetDefaultCorpusPath(t *testing.T) {
	corpusPath := os.Getenv("CORPUS_PATH")
	defer os.Setenv("CORPUS_PATH", corpusPath)
	os.Setenv("CORPUS_PATH", "")
	if getCorpusPath() != "corpus/" {
		t.Error(`default corpus path should be "corpus/"`)
	}
}
