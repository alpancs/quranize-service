package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alpancs/quranize/service"
	"github.com/go-chi/chi"
)

func Translate(w http.ResponseWriter, r *http.Request) {
	sura, errSura := strconv.Atoi(chi.URLParam(r, "sura"))
	aya, errAya := strconv.Atoi(chi.URLParam(r, "aya"))
	if errSura == nil && errAya == nil && validIndex(sura, aya) {
		json.NewEncoder(w).Encode(service.QuranTranslationID.Suras[sura-1].Ayas[aya-1].Text)
	} else {
		w.WriteHeader(400)
	}
}

func Tafsir(w http.ResponseWriter, r *http.Request) {
	sura, errSura := strconv.Atoi(chi.URLParam(r, "sura"))
	aya, errAya := strconv.Atoi(chi.URLParam(r, "aya"))
	if errSura == nil && errAya == nil && validIndex(sura, aya) {
		json.NewEncoder(w).Encode(service.QuranTafsirJalalayn.Suras[sura-1].Ayas[aya-1].Text)
	} else {
		w.WriteHeader(400)
	}
}

func validIndex(sura, aya int) bool {
	if sura < 1 || sura > len(service.QuranTranslationID.Suras) {
		return false
	}
	if aya < 1 || aya > len(service.QuranTranslationID.Suras[sura-1].Ayas) {
		return false
	}
	return true
}
