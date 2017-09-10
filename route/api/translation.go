package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alpancs/quranize/core"
	"github.com/go-chi/chi"
)

func Translation(w http.ResponseWriter, r *http.Request) {
	serve(w, r, core.QuranTranslationID)
}

func Tafsir(w http.ResponseWriter, r *http.Request) {
	serve(w, r, core.QuranTafsirQuraishShihab)
}

func Aya(w http.ResponseWriter, r *http.Request) {
	serve(w, r, core.QuranEnhanced)
}

func serve(w http.ResponseWriter, r *http.Request, quran core.Alquran) {
	sura, _ := strconv.Atoi(chi.URLParam(r, "sura"))
	aya, _ := strconv.Atoi(chi.URLParam(r, "aya"))
	if isValid(sura, aya) {
		json.NewEncoder(w).Encode(quran.Suras[sura-1].Ayas[aya-1].Text)
	} else {
		http.NotFound(w, r)
	}
}

func isValid(sura, aya int) bool {
	if sura < 1 || sura > len(core.QuranClean.Suras) {
		return false
	}
	if aya < 1 || aya > len(core.QuranClean.Suras[sura-1].Ayas) {
		return false
	}
	return true
}
