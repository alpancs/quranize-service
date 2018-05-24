package quranize

import (
	"strings"
)

var (
	hijaiyas       map[string][]string
	alphabetMaxLen int

	base = []string{""}
)

func Encode(text string) []string {
	var memo = make(map[string][]string)
	text = strings.Replace(text, " ", "", -1)
	text = strings.ToLower(text)
	results := []string{}
	for _, result := range quranize(text, memo) {
		if len(Locate(result)) > 0 {
			results = appendUniq(results, result)
		}
	}
	return results
}

func quranize(text string, memo map[string][]string) []string {
	if text == "" {
		return base
	}

	if cache, ok := memo[text]; ok {
		return cache
	}

	kalimas := []string{}
	l := len(text)
	for width := 1; width <= alphabetMaxLen && width <= l; width++ {
		if tails, ok := hijaiyas[text[l-width:]]; ok {
			heads := quranize(text[:l-width], memo)
			for _, combination := range combine(heads, tails) {
				if quran.exists(combination) {
					kalimas = appendUniq(kalimas, combination)
				}
			}
		}
	}

	memo[text] = kalimas
	return kalimas
}

func combine(heads, tails []string) []string {
	combinations := []string{}
	for _, head := range heads {
		for _, tail := range tails {
			combinations = append(combinations, head+tail)
			combinations = append(combinations, head+" "+tail)
			combinations = append(combinations, head+"ا"+tail)
			combinations = append(combinations, head+"ال"+tail)
			combinations = append(combinations, head+" ال"+tail)
			if tail == "و" {
				combinations = append(combinations, head+tail+"ا")
			}
		}
	}
	return combinations
}

func appendUniq(results []string, newResult string) []string {
	for _, result := range results {
		if result == newResult {
			return results
		}
	}
	return append(results, newResult)
}
