package quranize

var emptyLocation = []Location{}

func queryTree(harfs []rune, node *Node) []Location {
	for _, harf := range harfs {
		node = getChild(node.Children, harf)
		if node == nil {
			return emptyLocation
		}
	}
	return node.Locations
}

func Locate(kalima string) []Location {
	return queryTree([]rune(kalima), root)
}
