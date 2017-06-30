package route

import (
	"encoding/json"
	"net/http"

	"github.com/alpancs/quranize/service"
	"github.com/pressly/chi"
)

type (
	Location struct {
		Sura, Aya, Begin, End int
		SuraName, AyaText     string
	}
)

var (
	simpleHarfs = []rune{' ', 'ء', 'آ', 'أ', 'ؤ', 'إ', 'ئ', 'ا', 'ب', 'ة', 'ت', 'ث', 'ج', 'ح', 'خ', 'د', 'ذ', 'ر', 'ز', 'س', 'ش', 'ص', 'ض', 'ط', 'ظ', 'ع', 'غ', 'ف', 'ق', 'ك', 'ل', 'م', 'ن', 'ه', 'و', 'ى', 'ي'}
	quranClean  = &service.QuranClean
	quranMin    = &service.QuranMin
)

func Encode(w http.ResponseWriter, r *http.Request) {
	input := chi.URLParam(r, "input")
	json.NewEncoder(w).Encode(service.Encode(input))
}

func Locate(w http.ResponseWriter, r *http.Request) {
	input := chi.URLParam(r, "input")
	locations := []Location{}
	for _, loc := range service.Locate(input) {
		suraName := quranMin.Suras[loc.Sura].Name
		ayaText := quranMin.Suras[loc.Sura].Ayas[loc.Aya].Text
		begin := offset([]rune(ayaText), loc.Index)
		end := begin + len(input) - 2
		locations = append(locations, Location{loc.Sura, loc.Aya, begin, end, suraName, ayaText})
	}
	json.NewEncoder(w).Encode(locations)
}

func offset(runes []rune, n int) int {
	for i, harf := range runes {
		if n == 0 {
			return i
		}
		if isSimpleHarf(harf) {
			n--
		}
	}
	return len(runes)
}

func isSimpleHarf(harf rune) bool {
	begin, end := 0, len(simpleHarfs)-1
	for begin <= end {
		mid := (begin + end) / 2
		if harf > simpleHarfs[mid] {
			begin = mid + 1
		} else if harf < simpleHarfs[mid] {
			end = mid - 1
		} else {
			return true
		}
	}
	return false
}
