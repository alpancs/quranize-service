package quranize

import (
	"testing"

	"github.com/alpancs/quranize/corpus"
	"github.com/stretchr/testify/assert"
)

func TestLoadQuranBadXMLFormat(t *testing.T) {
	corpus.QuranSimpleCleanXML = ""
	assert.Panics(t, parseQuran)
}
