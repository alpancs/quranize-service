package quran

import (
	"encoding/xml"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/alpancs/quranize/quran/corpus"
)

type Transliteration struct {
	Hijaiyas map[string][]string
	MaxWidth int
}

var (
	QuranSimpleClean    Quran
	QuranEnhanced       Quran
	TranslationID       Quran
	TafsirQuraishShihab Quran

	transliteration = Transliteration{make(map[string][]string), 0}
)

func init() {
	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(5)
	go parseTransliterationAsync(&wg, corpus.ArabicToAlphabet, &transliteration)
	go parseAndBuildIndexAsync(&wg, corpus.QuranSimpleCleanXML, &QuranSimpleClean)
	go parseQuranAsync(&wg, corpus.QuranSimpleEnhancedXML, &QuranEnhanced)
	go parseQuranAsync(&wg, corpus.IDIndonesianXML, &TranslationID)
	go parseQuranAsync(&wg, corpus.IDMuntakhabXML, &TafsirQuraishShihab)
	wg.Wait()
	fmt.Println("service initialized in", time.Since(startTime))
}

func parseTransliterationAsync(wg *sync.WaitGroup, raw string, t *Transliteration) {
	parseTransliteration(raw, t)
	wg.Done()
}

func parseTransliteration(raw string, t *Transliteration) {
	trimmed := strings.TrimSpace(raw)
	for _, line := range strings.Split(trimmed, "\n") {
		components := strings.Split(line, " ")
		arabic := components[0]
		for _, alphabet := range components[1:] {
			t.Hijaiyas[alphabet] = append(t.Hijaiyas[alphabet], arabic)

			length := len(alphabet)
			ending := alphabet[length-1]
			if ending == 'a' || ending == 'i' || ending == 'o' || ending == 'u' {
				alphabet = alphabet[:length-1] + alphabet[:length-1] + alphabet[length-1:]
			} else {
				alphabet += alphabet
			}
			t.Hijaiyas[alphabet] = append(t.Hijaiyas[alphabet], arabic)
			length = len(alphabet)
			if length > t.MaxWidth {
				t.MaxWidth = length
			}
		}
	}
}

func parseQuranAsync(wg *sync.WaitGroup, raw string, q *Quran) {
	parseQuran(raw, q)
	wg.Done()
}

func parseAndBuildIndexAsync(wg *sync.WaitGroup, raw string, q *Quran) {
	parseQuran(raw, q)
	q.BuildIndex()
	wg.Done()
}

func parseQuran(raw string, q *Quran) {
	err := xml.Unmarshal([]byte(raw), q)
	if err != nil {
		panic(err)
	}
}
