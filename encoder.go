package quranize

import (
	"bytes"
	"strings"
)

const END = "$"

type (
	Location    struct{ Sura, Aya int }
	LocationSet map[Location]struct{}
	Children    map[rune]*Node

	Node struct {
		LocationSet
		Children
	}
)

var (
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

	root     = &Node{make(LocationSet), make(Children)}
	hijaiyas = make(map[string][]string)
	memo     = make(map[string][]string)
)

func queryTree(harfs []rune, node *Node) []Location {
	if node == nil {
		return []Location{}
	}
	if len(harfs) == 0 {
		locations := make([]Location, 0, len(node.LocationSet))
		for location := range node.LocationSet {
			locations = append(locations, location)
		}
		return locations
	}
	return queryTree(harfs[1:], node.Children[harfs[0]])
}

func inTree(harfs []rune, node *Node) bool {
	if node == nil {
		return false
	}
	if len(harfs) == 0 {
		return true
	}
	return inTree(harfs[1:], node.Children[harfs[0]])
}

func combine(heads, tails []string) []string {
	combinations := []string{}
	for _, head := range heads {
		for _, tail := range tails {
			if head == "" {
				combinations = append(combinations, tail)
			} else {
				combinations = append(combinations, head+tail)
				combinations = append(combinations, head+" "+tail)
				combinations = append(combinations, head+" ال"+tail)
				if head[len(head)-len("و"):] == "و" {
					combinations = append(combinations, head+"ا "+tail)
				}
			}
		}
	}
	return combinations
}

func quranize(text string) []string {
	if text == "" {
		return []string{""}
	}

	if cache, ok := memo[text]; ok {
		return cache
	}

	kalimas := []string{}
	l := len(text)
	for width := 1; width <= maxWidth && width <= l; width++ {
		if tails, ok := hijaiyas[text[l-width:]]; ok {
			heads := quranize(text[:l-width])
			for _, combination := range combine(heads, tails) {
				if inTree([]rune(combination), root) {
					kalimas = append(kalimas, combination)
				}
			}
		}
	}

	memo[text] = kalimas
	return kalimas
}

func removeDup(text string) string {
	buffer := bytes.NewBuffer(nil)
	for i := range text {
		if i == 0 || text[i-1] != text[i] {
			buffer.WriteByte(text[i])
		}
	}
	return buffer.String()
}

func Encode(text string) []string {
	text = strings.Replace(text, " ", "", -1)
	mixResults := append(quranize(text), quranize(removeDup(text))...)

	results := []string{}
	used := make(map[string]bool)
	for _, result := range mixResults {
		if !used[result] {
			used[result] = true
			results = append(results, result)
		}
	}
	return results
}
