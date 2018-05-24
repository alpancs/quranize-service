package quranize

import (
	"errors"
	"fmt"
)

type Quran struct {
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

type Node struct {
	locations []Location
	children  []Child
}

type Location struct{ Sura, Aya, SliceIndex int }

type Child struct {
	key   rune
	value *Node
}

var (
	quran Quran
	root  *Node

	zeroLocs = make([]Location, 0, 0)
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

// Get locations of s, matching whole word
func Locate(s string) []Location {
	harfs := []rune(s)
	node := root
	for _, harf := range harfs {
		node = getChild(node.children, harf)
		if node == nil {
			return zeroLocs
		}
	}
	return node.locations
}

// Check wether string s in Quran q or not
func (q Quran) exists(s string) bool {
	harfs := []rune(s)
	node := root
	for _, harf := range harfs {
		node = getChild(node.children, harf)
		if node == nil {
			return false
		}
	}
	return true
}

func buildIndex() {
	root = &Node{locations: zeroLocs}
	for _, sura := range quran.Suras {
		for _, aya := range sura.Ayas {
			indexAya([]rune(aya.Text), sura.Index, aya.Index, root)
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
		child := getChild(node.children, harf)
		if child == nil {
			child = &Node{}
			node.children = append(node.children, Child{harf, child})
		}
		node = child
		if i == len(harfs)-1 || harfs[i+1] == ' ' {
			node.locations = append(node.locations, location)
		}
	}
}

func getChild(children []Child, key rune) *Node {
	for _, child := range children {
		if child.key == key {
			return child.value
		}
	}
	return nil
}
