package quranize

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
)

type (
	Location struct{ Sura, Aya, Index int }
	Child    struct {
		Key   rune
		Value *Node
	}
	Node struct {
		Locations []Location
		Children  []Child
	}
)

var (
	maxWidth int

	root     = &Node{}
	hijaiyas = make(map[string][]string)

	Quran struct {
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
)

func getChild(children []Child, key rune) *Node {
	for _, child := range children {
		if child.Key == key {
			return child.Value
		}
	}
	return nil
}

func init() {
	loadHijaiyas()
	loadQuran()
	buildIndex()
}

func loadHijaiyas() {
	filePath := "corpus/arabic-to-alphabet"
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	trimmed := strings.TrimSpace(string(raw))
	for _, line := range strings.Split(trimmed, "\n") {
		components := strings.Split(line, " ")
		arabic := components[0]
		for _, alphabet := range components[1:] {
			hijaiyas[alphabet] = append(hijaiyas[alphabet], arabic)
			if maxWidth < len(alphabet) {
				maxWidth = len(alphabet)
			}
		}
	}
}

func loadQuran() {
	filePath := "corpus/quran-simple-clean.xml"
	raw, err := ioutil.ReadFile(filePath)
	if err == nil {
		err = xml.Unmarshal(raw, &Quran)
	}
	if err != nil {
		panic(err)
	}
}

func buildIndex() {
	for s, sura := range Quran.Suras {
		for a, aya := range sura.Ayas {
			indexAya([]rune(aya.Text), s, a)
		}
	}
}

func indexAya(harfs []rune, sura, aya int) {
	index := 0
	for index < len(harfs) {
		buildTree(harfs[index:], Location{sura, aya, index})
		for index < len(harfs) && harfs[index] != ' ' {
			index++
		}
		index++
	}
}

func buildTree(harfs []rune, location Location) {
	node := root
	for _, harf := range harfs {
		child := getChild(node.Children, harf)
		if child == nil {
			child = &Node{}
			node.Children = append(node.Children, Child{harf, child})
		}
		node = child
		node.Locations = appendUniqueLocation(node.Locations, location)
	}
}

func appendUniqueLocation(locations []Location, newLocation Location) []Location {
	for _, location := range locations {
		if newLocation == location {
			return locations
		}
	}
	return append(locations, newLocation)
}
