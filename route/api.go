package route

import (
	"encoding/json"
	"net/http"

	"github.com/alpancs/quranize/service"
	"github.com/pressly/chi"
)

type (
	CleanMin struct{ Clean, Min string }
	Location struct {
		Sura, Aya, Index  int
		SuraName, AyaText string
	}
)

var (
	simpleHarfs = []rune{' ', 'ء', 'آ', 'أ', 'ؤ', 'إ', 'ئ', 'ا', 'ب', 'ة', 'ت', 'ث', 'ج', 'ح', 'خ', 'د', 'ذ', 'ر', 'ز', 'س', 'ش', 'ص', 'ض', 'ط', 'ظ', 'ع', 'غ', 'ف', 'ق', 'ك', 'ل', 'م', 'ن', 'ه', 'و', 'ى', 'ي'}
	quranClean  = &service.QuranClean
	quranMin    = &service.QuranMin
)

func Encode(w http.ResponseWriter, r *http.Request) {
	input := chi.URLParam(r, "input")
	couples := []CleanMin{}
	for _, clean := range service.Encode(input) {
		couples = append(couples, CleanMin{clean, giveHarakah(clean)})
	}
	json.NewEncoder(w).Encode(couples)
}

func Locate(w http.ResponseWriter, r *http.Request) {
	input := chi.URLParam(r, "input")
	locations := []Location{}
	for _, loc := range service.Locate(input) {
		suraName := quranMin.Suras[loc.Sura].Name
		ayaText := quranMin.Suras[loc.Sura].Ayas[loc.Aya].Text
		index := offset([]rune(ayaText), loc.Index)
		locations = append(locations, Location{loc.Sura, loc.Aya, index, suraName, ayaText})
	}
	json.NewEncoder(w).Encode(locations)
}

func giveHarakah(kalima string) string {
	loc := service.Locate(kalima)[0]
	ayaMin := []rune(quranMin.Suras[loc.Sura].Ayas[loc.Aya].Text)
	begin := offset(ayaMin, loc.Index)
	end := begin + offset(ayaMin[begin:], len([]rune(kalima)))
	return string(ayaMin[begin : end+1])
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
