package main

import (
	"strings"
	"unicode"

	snowballeng "github.com/kljensen/snowball/english"
)

// step to create an inverted index
//

// raw text -> tokenizer -> filter -> tokens

func analyze(text string) []string {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopWordsFilter(tokens)
	return stemmerFilter(tokens)
}

func tokenize(text string) []string {
	// the strings.FieldsFunc splits the text if the rune
	// is not a letter or if it's a number
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}

	return r
}

func stopWordsFilter(tokens []string) []string {
	var stopWords = map[string]struct{}{
		"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
		"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
	}

	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopWords[token]; !ok {
			r = append(r, token)
		}
	}

	return r
}

// steaming reduces words into their base form
func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}

	return r
}
