package quran

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQuranBadXMLFormat(t *testing.T) {
	assert.Panics(t, func() { parseQuran("", nil) })
}
