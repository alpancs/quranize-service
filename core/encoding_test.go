package core

import "testing"

func testEncode(t *testing.T, input string, expected []string) {
	result := Encode(input)
	if !isStringListEqual(result, expected) {
		t.Errorf("expected: %v, result: %v", expected, result)
	}
}

func isStringListEqual(list1, list2 []string) bool {
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

func TestEncodeTajri(t *testing.T) {
	input := "tajri"
	output := []string{"تجري"}
	testEncode(t, input, output)
}

func TestEncodeAlhamdulillah(t *testing.T) {
	input := "alhamdulillah"
	output := []string{"الحمد لله"}
	testEncode(t, input, output)
}

func TestEncodeAlhamdulillahFull(t *testing.T) {
	input := "alhamdu lillahi robbil 'alamin"
	output := []string{"الحمد لله رب العالمين"}
	testEncode(t, input, output)
}

func TestEncodeBismillah(t *testing.T) {
	input := "bismillah"
	output := []string{"بسم الله", "بشماله"}
	testEncode(t, input, output)
}

func TestEncodeBismillahFull(t *testing.T) {
	input := "bismillah hirrohman nirrohim"
	output := []string{"بسم الله الرحمن الرحيم"}
	testEncode(t, input, output)
}

func TestEncodeWatasimu(t *testing.T) {
	input := "wa'tasimu"
	output := []string{"واعتصموا"}
	testEncode(t, input, output)
}

func TestEncodeWatasimuFull(t *testing.T) {
	input := "wa'tasimu bihablillah"
	output := []string{"واعتصموا بحبل الله"}
	testEncode(t, input, output)
}

func TestEncodeAlfatihah1(t *testing.T) {
	input := "bismilla hirrohma nirrohim"
	output := []string{"بسم الله الرحمن الرحيم"}
	testEncode(t, input, output)
}

func TestEncodeAlfatihah2(t *testing.T) {
	input := "alhamdu lillahi robbil 'alamin"
	output := []string{"الحمد لله رب العالمين"}
	testEncode(t, input, output)
}

func TestEncodeAlfatihah3(t *testing.T) {
	input := "arrohma nirrohim"
	output := []string{"الرحمن الرحيم"}
	testEncode(t, input, output)
}

func TestEncodeAlfatihah4(t *testing.T) {
	input := "maaliki yau middin"
	output := []string{"مالك يوم الدين"}
	testEncode(t, input, output)
}

func TestEncodeAlfatihah5(t *testing.T) {
	input := "iyya kanakbudu waiyya kanastain"
	output := []string{"إياك نعبد وإياك نستعين"}
	testEncode(t, input, output)
}

func TestEncodeAlfatihah6(t *testing.T) {
	input := "ihdinash shirothol mustaqim"
	output := []string{"اهدنا الصراط المستقيم"}
	testEncode(t, input, output)
}

func TestEncodeAlfatihah7(t *testing.T) {
	input := "shirotholladzina an'am ta'alaihim ghoiril maghdzu bi'alaihim waladh dhollin"
	output := []string{"صراط الذين أنعمت عليهم غير المغضوب عليهم ولا الضالين"}
	testEncode(t, input, output)
}

func TestShummunBukmun(t *testing.T) {
	input := "shummun bukmun"
	output := []string{"صم وبكم", "الصم البكم", "صم بكم"}
	testEncode(t, input, output)
}

func TestKahfi(t *testing.T) {
	input := "kahfi"
	output := []string{"الكهف"}
	testEncode(t, input, output)
}
