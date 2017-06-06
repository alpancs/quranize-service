package quranize

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

const END = "$"

type Location struct{ Sura, Aya int }

type None struct{}

type Node struct {
	Locations map[Location]None
	Children  map[rune]*Node
}

var (
	root     *Node
	hijaiyas = make(map[string][]string)
	maxWidth int

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

func init() {
	loadHijaiyas()
	loadQuran()
	preCompute()
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

func preCompute() {
	for s, sura := range Quran.Suras {
		for a, aya := range sura.Ayas {
			location := Location{s, a}
			buildIndex(aya.Text, location)
		}
	}
}

func buildIndex(text string, location Location) {
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
		node = &Node{make(map[Location]None), make(map[rune]*Node)}
	}
	node.Locations[location] = None{}
	node.Children[harf] = buildTree(text[width:], location, node.Children[harf])
	return node
}

func queryTree(text string, node *Node) []Location {
	if text == "" {
		locations := make([]Location, len(node.Locations))
		i := 0
		for location := range node.Locations {
			locations[i] = location
			i++
		}
		return locations
	}

	harf, width := utf8.DecodeRuneInString(text)
	if child, ok := node.Children[harf]; ok {
		return queryTree(text[width:], child)
	}
	return []Location{}
}

func combine(heads, tails []string) []string {
	combinations := []string{}
	for _, head := range heads {
		for _, tail := range tails {
			if tail == "" {
				combinations = append(combinations, head)
			} else {
				combinations = append(combinations, head+tail)
				combinations = append(combinations, head+" "+tail)
			}
		}
	}
	return combinations
}

func inQuran(text string) bool {
	return len(queryTree(text, root)) > 0
}

func quranize(text string) []string {
	if text == "" {
		return []string{""}
	}
	kalimas := []string{}
	l := len(text)
	for width := 1; width <= maxWidth && width <= l; width++ {
		if tails, ok := hijaiyas[text[l-width:]]; ok {
			heads := quranize(text[:l-width])
			for _, combination := range combine(heads, tails) {
				if inQuran(combination) {
					kalimas = append(kalimas, combination)
				}
			}
		}
	}
	return kalimas
}

func Encode(text string) []string {
	text = strings.Replace(text, " ", "", -1)
	return quranize(text)
}
