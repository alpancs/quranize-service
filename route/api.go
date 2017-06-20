package route

import (
	"encoding/json"
	"net/http"

	"github.com/alpancs/quranize/service"
	"github.com/pressly/chi"
)

type AyaLocation struct {
	Sura, Aya string
	Location  service.Location
}

var quran = &service.QuranSimple

func Encode(w http.ResponseWriter, r *http.Request) {
	input := chi.URLParam(r, "input")
	json.NewEncoder(w).Encode(service.Encode(input))
}

func Locate(w http.ResponseWriter, r *http.Request) {
	input := chi.URLParam(r, "input")
	locations := service.Locate(input)
	ret := make([]AyaLocation, len(locations))
	for i, location := range locations {
		sura := quran.Suras[location.Sura].Name
		aya := quran.Suras[location.Sura].Ayas[location.Aya].Text
		ret[i] = AyaLocation{sura, aya, location}
	}
	json.NewEncoder(w).Encode(ret)
}
