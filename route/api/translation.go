package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alpancs/quranize/core"
	"github.com/go-chi/chi"
)

func Translate(w http.ResponseWriter, r *http.Request) {
	sura, aya := getSuraAya(r)
	if isValid(sura, aya) {
		json.NewEncoder(w).Encode(core.QuranTranslationID.Suras[sura-1].Ayas[aya-1].Text)
	} else {
		w.WriteHeader(400)
	}
}

func Tafsir(w http.ResponseWriter, r *http.Request) {
	sura, aya := getSuraAya(r)
	if isValid(sura, aya) {
		json.NewEncoder(w).Encode(core.QuranTafsirQuraishShihab.Suras[sura-1].Ayas[aya-1].Text)
	} else {
		w.WriteHeader(400)
	}
}

func Aya(w http.ResponseWriter, r *http.Request) {
	sura, aya := getSuraAya(r)
	if isValid(sura, aya) {
		json.NewEncoder(w).Encode(core.QuranEnhanced.Suras[sura-1].Ayas[aya-1].Text)
	} else {
		w.WriteHeader(400)
	}
}

func getSuraAya(r *http.Request) (int, int) {
	sura, _ := strconv.Atoi(chi.URLParam(r, "sura"))
	aya, _ := strconv.Atoi(chi.URLParam(r, "aya"))
	return sura, aya
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
