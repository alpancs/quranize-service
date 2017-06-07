package quranize

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
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
		root = buildTree([]rune(text+END), location, root)
		start = strings.Index(text, " ") + 1
		if start == 0 {
			break
		}
	}
}

func buildTree(harfs []rune, location Location, node *Node) *Node {
	if len(harfs) == 0 {
		return node
	}

	if node == nil {
		node = &Node{make(LocationSet), make(Children)}
	}
	node.LocationSet[location] = none
	node.Children[harfs[0]] = buildTree(harfs[1:], location, node.Children[harfs[0]])
	return node
}
