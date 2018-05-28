package quran

import (
	"github.com/alpancs/quranize"
)

var (
	quran = quranize.NewDefaultQuranize()

	Encode = quran.Encode
	Locate = quran.Locate

	QuranSimpleEnhanced = quranize.NewQuranSimpleEnhanced()
	TranslationID       = quranize.NewIDIndonesian()
	TafsirID            = quranize.NewIDMuntakhab()
)
