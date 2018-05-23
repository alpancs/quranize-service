package quran

import (
	"errors"
	"fmt"
)

type Quran struct {
	Suras []Sura `xml:"sura"`
	root  *Node
}

type Sura struct {
	Index int    `xml:"index,attr"`
	Name  string `xml:"name,attr"`
	Ayas  []Aya  `xml:"aya"`
}

type Aya struct {
	Index     int    `xml:"index,attr"`
	Text      string `xml:"text,attr"`
	Bismillah string `xml:"bismillah,attr"`
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

var (
	emptyLocations = make([]Location, 0, 0)
)

// Get sura name from sura number (number starting from 1)
func (q Quran) GetSuraName(sura int) (string, error) {
	if !(1 <= sura && sura <= len(q.Suras)) {
		return "", errors.New(fmt.Sprintf("invalid sura number %d", sura))
	}
	return q.Suras[sura-1].Name, nil
}

// Get aya text from sura number and aya number (number starting from 1)
func (q Quran) GetAya(sura, aya int) (string, error) {
	if !(1 <= sura && sura <= len(q.Suras)) {
		return "", errors.New(fmt.Sprintf("invalid sura number %d", sura))
	}
	ayas := q.Suras[sura-1].Ayas
	if !(1 <= aya && aya <= len(ayas)) {
		return "", errors.New(fmt.Sprintf("invalid sura number %d and aya number %d", sura, aya))
	}
	return ayas[aya-1].Text, nil
}

// Build index of Quran q
func (q *Quran) BuildIndex() {
	q.root = &Node{Locations: emptyLocations}
	for _, sura := range q.Suras {
		for _, aya := range sura.Ayas {
			indexAya([]rune(aya.Text), sura.Index, aya.Index, q.root)
		}
	}
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

// Get locations of kalima in Quran q, matching whole word
func (q Quran) Locate(kalima string) []Location {
	harfs := []rune(kalima)
	node := q.root
	for _, harf := range harfs {
		node = getChild(node.Children, harf)
		if node == nil {
			return emptyLocations
		}
	}
	return node.Locations
}

// Check wether string s in Quran q or not
func (q Quran) Exists(s string) bool {
	harfs := []rune(s)
	node := q.root
	for _, harf := range harfs {
		node = getChild(node.Children, harf)
		if node == nil {
			return false
		}
	}
	return true
}
