package quran

import (
	"fmt"
	"sync"
	"time"

	"github.com/alpancs/quranize"
)

type Quran = quranize.Quran

var (
	q      quranize.Quranize
	Encode func(string) []string
	Locate func(string) []quranize.Location

	QuranSimpleEnhanced Quran
	TranslationID       Quran
	TafsirID            Quran

	wg sync.WaitGroup
)

func init() {
	startTime := time.Now()

	wg.Add(4)
	go initQuranize()
	go initQuranSimpleEnhanced()
	go initTranslationID()
	go initTafsirID()
	wg.Wait()

	fmt.Println("quran init time:", time.Since(startTime))
}

func initQuranize() {
	q = quranize.NewDefaultQuranize()
	Encode = q.Encode
	Locate = q.Locate
	wg.Done()
}

func initQuranSimpleEnhanced() {
	QuranSimpleEnhanced = quranize.NewQuranSimpleEnhanced()
	wg.Done()
}

func initTranslationID() {
	TranslationID = quranize.NewIDIndonesian()
	wg.Done()
}

func initTafsirID() {
	TafsirID = quranize.NewIDMuntakhab()
	wg.Done()
}
