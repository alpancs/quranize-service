package route

import (
	"encoding/json"
	"net/http"

	"github.com/alpancs/quranize/service"
	"github.com/pressly/chi"
)

type (
	CleanEnhanced struct{ Clean, Enhanced string }
	AyaLocation   struct {
		Sura, Aya string
		Location  service.Location
	}
)

var (
	simpleHarfs   = []rune{' ', 'ء', 'آ', 'أ', 'ؤ', 'إ', 'ئ', 'ا', 'ب', 'ة', 'ت', 'ث', 'ج', 'ح', 'خ', 'د', 'ذ', 'ر', 'ز', 'س', 'ش', 'ص', 'ض', 'ط', 'ظ', 'ع', 'غ', 'ف', 'ق', 'ك', 'ل', 'م', 'ن', 'ه', 'و', 'ى', 'ي'}
	quranClean    = &service.QuranClean
	quranEnhanced = &service.QuranEnhanced
)

func Encode(w http.ResponseWriter, r *http.Request) {
	input := chi.URLParam(r, "input")
	couples := []CleanEnhanced{}
	for _, clean := range service.Encode(input) {
		couples = append(couples, CleanEnhanced{clean, giveHarakah(clean)})
	}
	json.NewEncoder(w).Encode(couples)
}

func Locate(w http.ResponseWriter, r *http.Request) {
	input := chi.URLParam(r, "input")
	ret := []AyaLocation{}
	for _, loc := range service.Locate(input) {
		sura := quranClean.Suras[loc.Sura].Name
		aya := quranClean.Suras[loc.Sura].Ayas[loc.Aya].Text
		ret = append(ret, AyaLocation{sura, aya, loc})
	}
	json.NewEncoder(w).Encode(ret)
}

func giveHarakah(kalima string) string {
	loc := service.Locate(kalima)[0]
	ayaEnhanced := []rune(quranEnhanced.Suras[loc.Sura].Ayas[loc.Aya].Text)
	begin := offset(ayaEnhanced, loc.Index)
	end := begin + offset(ayaEnhanced[begin:], len([]rune(kalima)))
	return string(ayaEnhanced[begin : end+1])
}

func inSlice(target rune, slice []rune) bool {
	for _, e := range slice {
		if target == e {
			return true
		}
	}
	return false
}

func offset(runes []rune, n int) int {
	for i, harf := range runes {
		if n == 0 {
			return i
		}
		if inSlice(harf, simpleHarfs) {
			n--
		}
	}
	return -1
}
