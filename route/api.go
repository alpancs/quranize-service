package route

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/alpancs/quranize/service"
	"github.com/go-chi/chi"
)

type Location struct {
	Sura, Aya, Begin, End int
	SuraName, AyaText     string
}

func Encode(w http.ResponseWriter, r *http.Request) {
	keyword, _ := url.QueryUnescape(chi.URLParam(r, "keyword"))
	json.NewEncoder(w).Encode(service.Encode(keyword))
}

func Locate(w http.ResponseWriter, r *http.Request) {
	keyword := chi.URLParam(r, "keyword")
	locations := []Location{}
	for _, loc := range service.Locate(keyword) {
		suraName := service.QuranMin.Suras[loc.Sura].Name
		ayaText := service.QuranMin.Suras[loc.Sura].Ayas[loc.Aya].Text
		ayaTextRune := []rune(ayaText)
		begin := indexAfterSpaces(ayaTextRune, loc.SliceIndex)
		end := begin + indexAfterSpaces(ayaTextRune[begin:], strings.Count(keyword, " ")+1) - 1
		locations = append(locations, Location{loc.Sura, loc.Aya, begin, end, suraName, ayaText})
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
