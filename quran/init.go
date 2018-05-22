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

type Alquran struct {
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
	QuranClean               Alquran
	QuranEnhanced            Alquran
	QuranTranslationID       Alquran
	QuranTafsirQuraishShihab Alquran

	transliteration Transliteration
	emptyLocations  = make([]Location, 0, 0)
	corpusPath      = getCorpusPath()
)

func init() {
	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(5)
	go loadTransliterationAsync(&wg, "arabic-to-alphabet", &transliteration)
	go loadQuranAndIndexAsync(&wg, "quran-simple-clean.xml", &QuranClean)
	go loadQuranAsync(&wg, "quran-simple-enhanced.xml", &QuranEnhanced)
	go loadQuranAsync(&wg, "id.indonesian.xml", &QuranTranslationID)
	go loadQuranAsync(&wg, "id.muntakhab.xml", &QuranTafsirQuraishShihab)
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
	m := make(map[string][]string)
	maxWidth := 0
	raw, err := ioutil.ReadFile(corpusPath + fileName)
	if err != nil {
		panic(err)
	}
	trimmed := strings.TrimSpace(string(raw))
	for _, line := range strings.Split(trimmed, "\n") {
		components := strings.Split(line, " ")
		arabic := components[0]
		for _, alphabet := range components[1:] {
			m[alphabet] = append(m[alphabet], arabic)

			length := len(alphabet)
			ending := alphabet[length-1]
			if ending == 'a' || ending == 'i' || ending == 'o' || ending == 'u' {
				alphabet = alphabet[:length-1] + alphabet[:length-1] + alphabet[length-1:]
			} else {
				alphabet += alphabet
			}
			m[alphabet] = append(m[alphabet], arabic)
			length = len(alphabet)
			if length > maxWidth {
				maxWidth = length
			}
		}
	}
	t.Hijaiyas = m
	t.MaxWidth = maxWidth
}

func loadQuranAsync(wg *sync.WaitGroup, fileName string, quran *Alquran) {
	loadQuran(fileName, quran)
	wg.Done()
}

func loadQuranAndIndexAsync(wg *sync.WaitGroup, fileName string, quran *Alquran) {
	loadQuran(fileName, quran)
	quran.root = buildIndex(quran)
	wg.Done()
}

func loadQuran(fileName string, quran *Alquran) {
	raw, err := ioutil.ReadFile(corpusPath + fileName)
	if err != nil {
		panic(err)
	}
	err = xml.Unmarshal(raw, quran)
	if err != nil {
		panic(err)
	}
}

func buildIndex(quran *Alquran) *Node {
	node := &Node{Locations: emptyLocations}
	for s, sura := range QuranClean.Suras {
		for a, aya := range sura.Ayas {
			indexAya([]rune(aya.Text), s, a, node)
		}
	}
	return node
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
