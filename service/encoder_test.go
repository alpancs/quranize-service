package service

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

func TestEncodeAlfatihah1(t *testing.T) {
	input := "bismilla hirrohma nirrohim"
	output := []string{"بسم الله الرحمن الرحيم"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncodeAlfatihah2(t *testing.T) {
	input := "alhamdu lillahi robbil 'alamin"
	output := []string{"الحمد لله رب العالمين"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncodeAlfatihah3(t *testing.T) {
	input := "arrohma nirrohim"
	output := []string{"الرحمن الرحيم"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncodeAlfatihah4(t *testing.T) {
	input := "maliki yau middin"
	output := []string{"مالك يوم الدين"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncodeAlfatihah5(t *testing.T) {
	input := "iyya kanakbudu waiyya kanastain"
	output := []string{"إياك نعبد وإياك نستعين"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncodeAlfatihah6(t *testing.T) {
	input := "ihdinash shirothol mustaqim"
	output := []string{"اهدنا الصراط المستقيم"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncodeAlfatihah7(t *testing.T) {
	input := "shirotholladzina an'am ta'alaihim ghoiril maghdzu bi'alaihim waladh dhollin"
	output := []string{"صراط الذين أنعمت عليهم غير المغضوب عليهم ولا الضالين"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestShummunBukmun(t *testing.T) {
	input := "shummun bukmun"
	output := []string{"صم وبكم", "صم بكم"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}
