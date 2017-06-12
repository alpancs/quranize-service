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

func TestEncodeTajri(t *testing.T) {
	input := "tajri"
	output := []string{"تجري", "تأجر", "تجر", "تجار"}
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
	output := []string{"واعتصمو", "واعتصم"}
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

func TestLocate(t *testing.T) {
	input := "بسم الله الرحمن الرحيم"
	output := []Location{Location{0, 0, 0, 22}, Location{26, 29, 9, 31}}
	if !properlyLocated(input, output) {
		t.Error(Locate(input))
	}
}
