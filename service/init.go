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

type Location struct{ Sura, Aya, Index int }

type Child struct {
	Key   rune
	Value *Node
}

type Node struct {
	Locations []Location
	Children  []Child
}

var (
	QuranClean, QuranMin Alquran

	maxWidth int
	root     *Node
	hijaiyas map[string][]string
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
	loadHijaiyas("corpus/arabic-to-alphabet")
	loadQuran("corpus/quran-simple-clean.xml", &QuranClean)
	loadQuran("corpus/quran-simple-min.xml", &QuranMin)
	root = buildIndex(&QuranClean)
}

func loadHijaiyas(filePath string) {
	hijaiyas = make(map[string][]string)
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
			if length > maxWidth {
				maxWidth = length
			}
		}
	}
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
	index := 0
	for index < len(harfs) {
		buildTree(harfs[index:], Location{sura, aya, index}, node)
		for index < len(harfs) && harfs[index] != ' ' {
			index++
		}
		index++
	}
}

func buildTree(harfs []rune, location Location, node *Node) {
	for _, harf := range harfs {
		child := getChild(node.Children, harf)
		if child == nil {
			child = &Node{}
			node.Children = append(node.Children, Child{harf, child})
		}
		node = child
		node.Locations = append(node.Locations, location)
	}
}
