package quranize

import (
	"bytes"
	"strings"
)

type (
	Location struct{ Sura, Aya, Begin, End int }
	Children map[rune]*Node

	Node struct {
		Locations []Location
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

	root          = &Node{Children: make(Children)}
	hijaiyas      = make(map[string][]string)
	memo          = make(map[string][]string)
	emptyLocation = []Location{}
)

func queryTree(harfs []rune, node *Node) []Location {
	for _, harf := range harfs {
		if node.Children[harf] == nil {
			return emptyLocation
		}
		node = node.Children[harf]
	}
	return node.Locations
}

func inTree(harfs []rune, node *Node) bool {
	for _, harf := range harfs {
		if node.Children[harf] == nil {
			return false
		}
		node = node.Children[harf]
	}
	return true
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
				combinations = append(combinations, head+"ال"+tail)
				combinations = append(combinations, head+" ال"+tail)
				combinations = append(combinations, head+"ا"+tail)
				combinations = append(combinations, head+"ا "+tail)
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

func removeDoubleChar(text string) string {
	buffer := bytes.NewBuffer(nil)
	for i := range text {
		if i == 0 || text[i-1] != text[i] {
			buffer.WriteByte(text[i])
		}
	}
	return buffer.String()
}

func unique(results []string) []string {
	uniqueResult := []string{}
	inSlice := make(map[string]bool)
	for _, result := range results {
		if !inSlice[result] {
			inSlice[result] = true
			uniqueResult = append(uniqueResult, result)
		}
	}
	return uniqueResult
}

func Encode(text string) []string {
	text = strings.Replace(text, " ", "", -1)
	results1 := quranize(text)
	results2 := quranize(removeDoubleChar(text))
	return unique(append(results1, results2...))
}

func Locate(kalima string) []Location {
	return queryTree([]rune(kalima), root)
}
