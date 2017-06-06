package quranize

import "testing"

func TestEncodeTajri(t *testing.T) {
	input := "tajri"
	output := []string{"تجري", "تجر"}

	kalimas := Encode(input)

	if len(kalimas) != len(output) {
		t.Fail()
		return
	}
	for i, kalima := range kalimas {
		if output[i] != kalima {
			t.Fail()
		}
	}
}

func TestEncodeAlhamdulilah(t *testing.T) {
	input := "alhamdulilah"
	output := []string{"الحمد لله"}

	kalimas := Encode(input)

	if len(kalimas) != len(output) {
		t.Fail()
		return
	}
	for i, kalima := range kalimas {
		if output[i] != kalima {
			t.Fail()
		}
	}
}
