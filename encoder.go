package quranize

import (
	"strings"
	"unicode/utf8"
)

const END = "$"

type Location struct{ Sura, Aya int }
type LocationSet map[Location]struct{}
type Children map[rune]*Node

type Node struct {
	LocationSet
	Children
}

var (
	root     *Node
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

	hijaiyas = make(map[string][]string)
)

func queryTree(text string, node *Node) []Location {
	if text == "" {
		locations := make([]Location, 0, len(node.LocationSet))
		for location := range node.LocationSet {
			locations = append(locations, location)
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
