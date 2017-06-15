package quranize

import "testing"

func properlyEncoded(input string, output []string) bool {
	kalimas := Encode(input)
	if len(kalimas) != len(output) {
		return false
	}
	for i, kalima := range kalimas {
		if output[i] != kalima {
			return false
		}
	}
	return true
}

func TestEncodeTajri(t *testing.T) {
	input := "tajri"
	output := []string{"تجري"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncodeAlhamdulillah(t *testing.T) {
	input := "alhamdulillah"
	output := []string{"الحمد لله"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncodeAlhamdulillahFull(t *testing.T) {
	input := "alhamdu lillahi robbil 'alamin"
	output := []string{"الحمد لله رب العالمين"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncodeBismillah(t *testing.T) {
	input := "bismillah"
	output := []string{"بسم الله", "بشماله"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncodeBismillahFull(t *testing.T) {
	input := "bismillah hirrohman nirrohim"
	output := []string{"بسم الله الرحمن الرحيم"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncodeWatasimu(t *testing.T) {
	input := "wa'tasimu"
	output := []string{"واعتصموا"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncodeWatasimuFull(t *testing.T) {
	input := "wa'tasimu bihablillah"
	output := []string{"واعتصموا بحبل الله"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}
