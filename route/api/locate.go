package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alpancs/quranize-service/quran"
)

type Location struct {
	SuraNumber     int    `json:"suraNumber"`
	SuraName       string `json:"suraName"`
	AyaNumber      int    `json:"ayaNumber"`
	AyaText        string `json:"ayaText"`
	BeginHighlight int    `json:"beginHighlight"`
	EndHighlight   int    `json:"endHighlight"`
}

func Locate(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	locations := []Location{}
	for _, location := range quran.Locate(keyword) {
		suraName, _ := quran.QuranSimpleEnhanced.GetSuraName(location.Sura)
		ayaText, _ := quran.QuranSimpleEnhanced.GetAya(location.Sura, location.Aya)
		ayaTextRune := []rune(ayaText)
		begin := indexAfterSpaces(ayaTextRune, location.WordIndex)
		end := begin + indexAfterSpaces(ayaTextRune[begin:], strings.Count(keyword, " ")+1) - 1
		locations = append(locations, Location{location.Sura, suraName, location.Aya, ayaText, begin, end})
	}
	json.NewEncoder(w).Encode(locations)
}

func indexAfterSpaces(text []rune, remain int) int {
	for i, r := range text {
		if remain == 0 {
			return i
		}
		if r == ' ' {
			remain--
		}
	}
	return len(text) + 1
}
