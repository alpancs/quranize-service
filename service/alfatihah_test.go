package service

import "testing"

func TestEncode1(t *testing.T) {
	input := "bismilla hirrohma nirrohim"
	output := []string{"بسم الله الرحمن الرحيم"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncode2(t *testing.T) {
	input := "alhamdu lillahi robbil 'alamin"
	output := []string{"الحمد لله رب العالمين"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncode3(t *testing.T) {
	input := "arrohma nirrohim"
	output := []string{"الرحمن الرحيم"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncode4(t *testing.T) {
	input := "maliki yau middin"
	output := []string{"مالك يوم الدين"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncode5(t *testing.T) {
	input := "iyya kanakbudu waiyya kanastain"
	output := []string{"إياك نعبد وإياك نستعين"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncode6(t *testing.T) {
	input := "ihdinash shirothol mustaqim"
	output := []string{"اهدنا الصراط المستقيم"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}

func TestEncode7(t *testing.T) {
	input := "shirotholladzina an'am ta'alaihim ghoiril maghdzu bi'alaihim waladh dhollin"
	output := []string{"صراط الذين أنعمت عليهم غير المغضوب عليهم ولا الضالين"}
	if !properlyEncoded(input, output) {
		t.Error(Encode(input))
	}
}
