package service

import "testing"

func properlyLocated(input string, output []Location) bool {
	locations := Locate(input)
	if len(locations) != len(output) {
		return false
	}
	for i, location := range locations {
		if output[i] != location {
			return false
		}
	}
	return true
}

func TestLocate(t *testing.T) {
	input := "بسم الله الرحمن الرحيم"
	output := []Location{Location{0, 0, 0}, Location{26, 29, 19}}
	if !properlyLocated(input, output) {
		t.Error(Locate(input))
	}
}
