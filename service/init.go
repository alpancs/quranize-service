package service

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
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
}

type Location struct{ Sura, Aya, SliceIndex int }

type Child struct {
	Key   rune
	Value *Node
}

type Node struct {
	Locations []Location
	Children  []Child
}

var (
	QuranClean          Alquran
	QuranEnhanced       Alquran
	QuranTranslationID  Alquran
	QuranTafsirJalalayn Alquran

	hijaiyas map[string][]string
	maxWidth int
	root     *Node
)

func init() {
	hijaiyas = loadTransliteration("corpus/arabic-to-alphabet")
	loadQuran("corpus/quran-simple-clean.xml", &QuranClean)
	root = buildIndex(&QuranClean)
	loadQuran("corpus/quran-simple-enhanced.xml", &QuranEnhanced)
	loadQuran("corpus/id.indonesian.xml", &QuranTranslationID)
	loadQuran("corpus/id.jalalayn.xml", &QuranTafsirJalalayn)
}

func loadTransliteration(filePath string) map[string][]string {
	dictionary := make(map[string][]string)
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		raw, err = ioutil.ReadFile("../" + filePath)
		if err != nil {
			panic(err)
		}
	}
	trimmed := strings.TrimSpace(string(raw))
	for _, line := range strings.Split(trimmed, "\n") {
		components := strings.Split(line, " ")
		arabic := components[0]
		for _, alphabet := range components[1:] {
			dictionary[alphabet] = append(dictionary[alphabet], arabic)

			length := len(alphabet)
			ending := alphabet[length-1]
			if ending == 'a' || ending == 'i' || ending == 'o' || ending == 'u' {
				alphabet = alphabet[:length-1] + alphabet[:length-1] + alphabet[length-1:]
			} else {
				alphabet += alphabet
			}
			dictionary[alphabet] = append(dictionary[alphabet], arabic)
			length = len(alphabet)
			if length > maxWidth {
				maxWidth = length
			}
		}
	}
	return dictionary
}

func loadQuran(filePath string, quran *Alquran) {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		raw, err = ioutil.ReadFile("../" + filePath)
		if err != nil {
			panic(err)
		}
	}
	err = xml.Unmarshal(raw, quran)
	if err != nil {
		panic(err)
	}
}

func buildIndex(quran *Alquran) *Node {
	node := &Node{}
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
