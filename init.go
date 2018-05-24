package quranize

import (
	"encoding/xml"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/alpancs/quranize/corpus"
)

var (
	QuranSimpleClean Quran
	hijaiyas         map[string][]string
	alphabetMaxLen   int
)

func init() {
	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(2)
	go parseQuranAsync(&wg)
	go parseDictAsync(&wg)
	wg.Wait()
	fmt.Println("package quranize initialized in", time.Since(startTime))
}

func parseQuranAsync(wg *sync.WaitGroup) {
	parseQuran()
	QuranSimpleClean.buildIndex()
	wg.Done()
}

func parseQuran() {
	err := xml.Unmarshal([]byte(corpus.QuranSimpleCleanXML), &QuranSimpleClean)
	if err != nil {
		panic(err)
	}
}

func parseDictAsync(wg *sync.WaitGroup) {
	parseDict()
	wg.Done()
}

func parseDict() {
	hijaiyas = make(map[string][]string)
	trimmed := strings.TrimSpace(corpus.ArabicToAlphabet)
	for _, line := range strings.Split(trimmed, "\n") {
		components := strings.Split(line, " ")
		arabic := components[0]
		for _, alphabet := range components[1:] {
			hijaiyas[alphabet] = append(hijaiyas[alphabet], arabic)

			length := len(alphabet)
			ending := alphabet[length-1]
			if ending == 'a' || ending == 'i' || ending == 'o' || ending == 'u' {
				alphabet = alphabet[:length-1] + alphabet[:length-1] + alphabet[length-1:]
			} else {
				alphabet += alphabet
			}
			hijaiyas[alphabet] = append(hijaiyas[alphabet], arabic)
			length = len(alphabet)
			if length > alphabetMaxLen {
				alphabetMaxLen = length
			}
		}
	}
}
