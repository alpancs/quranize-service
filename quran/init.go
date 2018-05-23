package quran

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type Quran struct {
	Suras []struct {
		Index int    `xml:"index,attr"`
		Name  string `xml:"name,attr"`
		Ayas  []struct {
			Index     int    `xml:"index,attr"`
			Text      string `xml:"text,attr"`
			Bismillah string `xml:"bismillah,attr"`
		} `xml:"aya"`
	} `xml:"sura"`
	root *Node
}

type Node struct {
	Locations []Location
	Children  []Child
}

type Location struct{ Sura, Aya, SliceIndex int }

type Child struct {
	Key   rune
	Value *Node
}

type Transliteration struct {
	Hijaiyas map[string][]string
	MaxWidth int
}

var (
	QuranClean          Quran
	QuranEnhanced       Quran
	TranslationID       Quran
	TafsirQuraishShihab Quran

	transliteration = Transliteration{make(map[string][]string), 0}
	emptyLocations  = make([]Location, 0, 0)
	corpusPath      = getCorpusPath()
)

func init() {
	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(5)
	go loadTransliterationAsync(&wg, "arabic-to-alphabet", &transliteration)
	go loadAndIndexQuranAsync(&wg, "quran-simple-clean.xml", &QuranClean)
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

// Get aya text from sura number and aya number (number starting from 1)
func (q Quran) GetAya(sura, aya int) (string, error) {
	if !(1 <= sura && sura <= len(q.Suras)) {
		return "", errors.New(fmt.Sprintf("invalid sura number %d", sura))
	}
	ayas := q.Suras[sura-1].Ayas
	if !(1 <= aya && aya <= len(ayas)) {
		return "", errors.New(fmt.Sprintf("invalid sura number %d and aya number %d", sura, aya))
	}
	return ayas[aya-1].Text, nil
}

// Build index of Quran q
func (q *Quran) BuildIndex() {
	q.root = &Node{Locations: emptyLocations}
	for s, sura := range QuranClean.Suras {
		for a, aya := range sura.Ayas {
			indexAya([]rune(aya.Text), s, a, q.root)
		}
	}
}

func indexAya(harfs []rune, sura, aya int, node *Node) {
	sliceIndex := 0
	for i := range harfs {
		if i == 0 || harfs[i-1] == ' ' {
			buildTree(harfs[i:], Location{sura, aya, sliceIndex}, node)
			sliceIndex++
		}
	}
}

func buildTree(harfs []rune, location Location, node *Node) {
	for i, harf := range harfs {
		child := getChild(node.Children, harf)
		if child == nil {
			child = &Node{}
			node.Children = append(node.Children, Child{harf, child})
		}
		node = child
		if i == len(harfs)-1 || harfs[i+1] == ' ' {
			node.Locations = append(node.Locations, location)
		}
	}
}

func getChild(children []Child, key rune) *Node {
	for _, child := range children {
		if child.Key == key {
			return child.Value
		}
	}
	return nil
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

func loadAndIndexQuranAsync(wg *sync.WaitGroup, fileName string, q *Quran) {
	loadQuran(fileName, q)
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
