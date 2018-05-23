package quran

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type Transliteration struct {
	Hijaiyas map[string][]string
	MaxWidth int
}

var (
	QuranEnhanced       Quran
	TranslationID       Quran
	TafsirQuraishShihab Quran

	transliteration = Transliteration{make(map[string][]string), 0}
	corpusPath      = getCorpusPath()
)

func init() {
	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(5)
	go loadTransliterationAsync(&wg, "arabic-to-alphabet", &transliteration)
	go buildIndexAsync(&wg, &QuranSimpleClean)
	go loadQuranAsync(&wg, "quran-simple-enhanced.xml", &QuranEnhanced)
	go loadQuranAsync(&wg, "id.indonesian.xml", &TranslationID)
	go loadQuranAsync(&wg, "id.muntakhab.xml", &TafsirQuraishShihab)
	wg.Wait()
	fmt.Println("service initialized in", time.Since(startTime))
}

func getCorpusPath() string {
	if path := os.Getenv("CORPUS_PATH"); path != "" {
		return path
	}
	return "corpus/"
}

func loadTransliterationAsync(wg *sync.WaitGroup, fileName string, t *Transliteration) {
	loadTransliteration(fileName, t)
	wg.Done()
}

func loadTransliteration(fileName string, t *Transliteration) {
	raw, err := ioutil.ReadFile(corpusPath + fileName)
	if err != nil {
		panic(err)
	}
	trimmed := strings.TrimSpace(string(raw))
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

func loadQuranAsync(wg *sync.WaitGroup, fileName string, q *Quran) {
	loadQuran(fileName, q)
	wg.Done()
}

func buildIndexAsync(wg *sync.WaitGroup, q *Quran) {
	q.BuildIndex()
	wg.Done()
}

func loadQuran(fileName string, q *Quran) {
	raw, err := ioutil.ReadFile(corpusPath + fileName)
	if err != nil {
		panic(err)
	}
	err = xml.Unmarshal(raw, q)
	if err != nil {
		panic(err)
	}
}
