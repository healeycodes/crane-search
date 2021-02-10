// Simple Full-Text Search engine
// Copied/Modified from https://github.com/akrylysov/simplefts
// LICENSE: https://github.com/akrylysov/simplefts/blob/master/LICENSE

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	idx := make(Index)

	assert.Nil(t, idx.Search("foo"))
	assert.Nil(t, idx.Search("donut"))

	idx.Add([]document{{ID: 1, Text: "A donut on a glass plate. Only the donuts."}})
	assert.Nil(t, idx.Search("a"))
	assert.Equal(t, idx.Search("donut"), []int{1})
	assert.Equal(t, idx.Search("DoNuts"), []int{1})
	assert.Equal(t, idx.Search("glass"), []int{1})

	idx.Add([]document{{ID: 2, Text: "donut is a donut"}})
	assert.Nil(t, idx.Search("a"))
	assert.Equal(t, idx.Search("donut"), []int{1, 2})
	assert.Equal(t, idx.Search("DoNuts"), []int{1, 2})
	assert.Equal(t, idx.Search("glass"), []int{1})
}

func TestTokenizer(t *testing.T) {
	testCases := []struct {
		text   string
		tokens []string
	}{
		{
			text:   "",
			tokens: []string{},
		},
		{
			text:   "a",
			tokens: []string{"a"},
		},
		{
			text:   "small wild,cat!",
			tokens: []string{"small", "wild", "cat"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.text, func(st *testing.T) {
			assert.EqualValues(st, tc.tokens, Tokenize(tc.text))
		})
	}
}

func TestLowercaseFilter(t *testing.T) {
	var (
		in  = []string{"Cat", "DOG", "fish"}
		out = []string{"cat", "dog", "fish"}
	)
	assert.Equal(t, out, LowercaseFilter(in))
}

func TestStopwordFilter(t *testing.T) {
	var (
		in  = []string{"i", "am", "the", "cat"}
		out = []string{"am", "cat"}
	)
	assert.Equal(t, out, StopwordFilter(in))
}

func TestStemmerFilter(t *testing.T) {
	var (
		in  = []string{"cat", "cats", "fish", "fishing", "fished", "airline"}
		out = []string{"cat", "cat", "fish", "fish", "fish", "airlin"}
	)
	assert.Equal(t, out, StemmerFilter(in))
}
