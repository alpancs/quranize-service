package quran

import "testing"

func testLocate(t *testing.T, input string, expected []Location) {
	result := Locate(input)
	if !isLocationListEqual(result, expected) {
		t.Errorf("input: %v, expected: %v, result: %v", input, expected, result)
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

func TestLocateEmptyString(t *testing.T) {
	input := ""
	output := []Location{}
	testLocate(t, input, output)
}

func TestLocateNonAlquran(t *testing.T) {
	input := "alfan"
	output := []Location{}
	testLocate(t, input, output)
}

func TestLocateAlquran(t *testing.T) {
	input := "بسم الله الرحمن الرحيم"
	output := []Location{Location{0, 0, 0}, Location{26, 29, 4}}
	testLocate(t, input, output)
}
