package route

import (
	"encoding/json"
	"net/http"

	"github.com/alpancs/quranize/service"
	"github.com/pressly/chi"
)

type (
	CleanEnhanced struct{ Clean, Enhanced string }
	Location      struct {
		Sura, Aya, Index  int
		SuraName, AyaText string
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
	locations := []Location{}
	for _, loc := range service.Locate(input) {
		suraName := quranEnhanced.Suras[loc.Sura].Name
		ayaText := quranEnhanced.Suras[loc.Sura].Ayas[loc.Aya].Text
		index := offset([]rune(ayaText), loc.Index)
		locations = append(locations, Location{loc.Sura, loc.Aya, index, suraName, ayaText})
	}
	json.NewEncoder(w).Encode(locations)
}

func giveHarakah(kalima string) string {
	loc := service.Locate(kalima)[0]
	ayaEnhanced := []rune(quranEnhanced.Suras[loc.Sura].Ayas[loc.Aya].Text)
	begin := offset(ayaEnhanced, loc.Index)
	end := begin + offset(ayaEnhanced[begin:], len([]rune(kalima)))
	return string(ayaEnhanced[begin : end+1])
}

func offset(runes []rune, n int) int {
	for i, harf := range runes {
		if n == 0 {
			return i
		}
		if inSlice(harf, simpleHarfs, 0, len(simpleHarfs)-1) {
			n--
		}
	}
	return -1
}

func inSlice(target rune, slice []rune, begin, end int) bool {
	for begin <= end {
		mid := (begin + end) / 2
		if target > slice[mid] {
			begin = mid + 1
		} else if target < slice[mid] {
			end = mid - 1
		} else {
			return true
		}
	}
	return false
}
