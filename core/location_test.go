package core

import "testing"

func testLocate(t *testing.T, input string, expected []Location) {
	result := Locate(input)
	if !isLocationListEqual(result, expected) {
		t.Errorf("expected: %v, result: %v", expected, result)
	}
}

func isLocationListEqual(list1, list2 []Location) bool {
	if len(list1) != len(list2) {
		return false
	}
	for i := range list1 {
		if list1[i] != list2[i] {
			return false
		}
	}
	return true
}

func TestLocateBismillahFull(t *testing.T) {
	input := "بسم الله الرحمن الرحيم"
	output := []Location{Location{0, 0, 0}, Location{26, 29, 4}}
	testLocate(t, input, output)
}
