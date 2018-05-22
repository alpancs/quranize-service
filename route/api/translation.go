package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alpancs/quranize/quran"
	"github.com/go-chi/chi"
)

func Translation(w http.ResponseWriter, r *http.Request) {
	serve(w, r, quran.QuranTranslationID)
}

func Tafsir(w http.ResponseWriter, r *http.Request) {
	serve(w, r, quran.QuranTafsirQuraishShihab)
}

func Aya(w http.ResponseWriter, r *http.Request) {
	serve(w, r, quran.QuranEnhanced)
}

func serve(w http.ResponseWriter, r *http.Request, q quran.Quran) {
	sura, _ := strconv.Atoi(chi.URLParam(r, "sura"))
	aya, _ := strconv.Atoi(chi.URLParam(r, "aya"))
	text, err := q.GetAya(sura, aya)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(text)
}
