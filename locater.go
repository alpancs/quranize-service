package quranize

var emptyLocation = []Location{}

func queryTree(harfs []rune, node *Node) []Location {
	for _, harf := range harfs {
		if node.Children[harf] == nil {
			return emptyLocation
		}
		node = node.Children[harf]
	}
	return node.Locations
}

func Locate(kalima string) []Location {
	return queryTree([]rune(kalima), root)
}
