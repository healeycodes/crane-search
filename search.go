// Uses code from Simple Full-Text Search engine
// Copied/Modified from https://artem.krylysov.com/blog/2020/07/28/lets-build-a-full-text-search-engine/
// LICENSE: https://github.com/akrylysov/simplefts/blob/master/LICENSE

package search

import (
	"strings"
	"unicode"

	snowballeng "github.com/kljensen/snowball/english"
)

// Document represents a text file
type Document struct {
	Title string
	URL   string
	Text  string
	ID    int
}

// Result is a search result item
type Result struct {
	Title string
	URL   string
	ID    int
}

// Index is an inverted Index. It maps tokens to document IDs.
type Index map[string][]int

// Store contains results and their index
type Store struct {
	Index   Index
	Results []Result
}

// Add adds documents to the index.
func (index Index) Add(docs []Document) {
	for _, doc := range docs {
		for _, token := range Analyze(doc.Text) {
			ids := index[token]
			if ids != nil && ids[len(ids)-1] == doc.ID {
				// Don't add same ID twice.
				continue
			}
			index[token] = append(ids, doc.ID)
		}
	}
}

// Intersection returns the set Intersection between a and b.
// a and b have to be sorted in ascending order and contain no duplicates.
func Intersection(a []int, b []int) []int {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	r := make([]int, 0, maxLen)
	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			r = append(r, a[i])
			i++
			j++
		}
	}
	return r
}

// Search queries the index for the given text.
func (index Index) Search(text string) []int {
	var r []int
	for _, token := range Analyze(text) {
		if ids, ok := index[token]; ok {
			if r == nil {
				r = ids
			} else {
				r = Intersection(r, ids)
			}
		} else {
			// Token doesn't exist.
			return nil
		}
	}
	return r
}

// Tokenize returns a slice of tokens for the given text.
func Tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// Analyze analyzes the text and returns a slice of tokens.
func Analyze(text string) []string {
	tokens := Tokenize(text)
	tokens = LowercaseFilter(tokens)
	tokens = StopwordFilter(tokens)
	tokens = StemmerFilter(tokens)
	return tokens
}

// LowercaseFilter returns a slice of tokens normalized to lower case.
func LowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

// StopwordFilter returns a slice of tokens with stop words removed.
func StopwordFilter(tokens []string) []string {
	var stopwords = map[string]struct{}{
		"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
		"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
	}
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopwords[token]; !ok {
			r = append(r, token)
		}
	}
	return r
}

// StemmerFilter returns a slice of stemmed tokens.
func StemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}
	return r
}
