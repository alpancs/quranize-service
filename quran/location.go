package quran

func queryTree(harfs []rune) []Location {
	node := root
	for _, harf := range harfs {
		node = getChild(node.Children, harf)
		if node == nil {
			return emptyLocations
		}
	}
	return node.Locations
}

func Locate(kalima string) []Location {
	return queryTree([]rune(kalima))
}
