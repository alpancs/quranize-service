package core

import "testing"

type TestCase struct {
	input  string
	output []string
}

func testEncode(t *testing.T, input string, expected []string) {
	result := Encode(input)
	if !isStringListEqual(result, expected) {
		t.Errorf("input: %v, expected: %v, result: %v", input, expected, result)
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

func TestEncodeEmptyString(t *testing.T) {
	input := ""
	output := []string{}
	testEncode(t, input, output)
}

func TestEncodeNonAlquran(t *testing.T) {
	input := "alfan nur fauzan"
	output := []string{}
	testEncode(t, input, output)
}

func TestEncodeAlquran(t *testing.T) {
	testCases := []TestCase{
		TestCase{"tajri", []string{"تجري"}},
		TestCase{"alhamdulillah", []string{"الحمد لله"}},
		TestCase{"bismillah", []string{"بسم الله", "بشماله"}},
		TestCase{"wa'tasimu", []string{"واعتصموا"}},
		TestCase{"wa'tasimu bihablillah", []string{"واعتصموا بحبل الله"}},
		TestCase{"shummun bukmun", []string{"صم وبكم", "صم بكم", "الصم البكم"}},
		TestCase{"kahfi", []string{"الكهف"}},
		TestCase{"wabasyiris sobirin", []string{"وبشر الصابرين"}},
		TestCase{"bissobri wassolah", []string{"بالصبر والصلاة"}},

		TestCase{"bismillah hirrohman nirrohim", []string{"بسم الله الرحمن الرحيم"}},
		TestCase{"alhamdu lillahi robbil 'alamin", []string{"الحمد لله رب العالمين"}},
		TestCase{"arrohma nirrohim", []string{"الرحمن الرحيم"}},
		TestCase{"maaliki yau middin", []string{"مالك يوم الدين"}},
		TestCase{"iyya kanakbudu waiyya kanastain", []string{"إياك نعبد وإياك نستعين"}},
		TestCase{"ihdinash shirothol mustaqim", []string{"اهدنا الصراط المستقيم"}},
		TestCase{"shirotholladzina an'am ta'alaihim ghoiril maghdzu bi'alaihim waladh dhollin", []string{"صراط الذين أنعمت عليهم غير المغضوب عليهم ولا الضالين"}},
	}
	for _, testCase := range testCases {
		testEncode(t, testCase.input, testCase.output)
	}
}
