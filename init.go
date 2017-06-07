package quranize

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

var none struct{}

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
			location := Location{s, a}
			indexAya(aya.Text, location)
		}
	}
}

func indexAya(text string, location Location) {
	start := 0
	for {
		text = text[start:]
		root = buildTree(text+END, location, root)
		start = strings.Index(text, " ") + 1
		if start == 0 {
			break
		}
	}
}

func buildTree(text string, location Location, node *Node) *Node {
	if text == "" {
		return node
	}

	harf, width := utf8.DecodeRuneInString(text)
	if node == nil {
		node = &Node{make(LocationSet), make(Children)}
	}
	node.LocationSet[location] = none
	node.Children[harf] = buildTree(text[width:], location, node.Children[harf])
	return node
}
