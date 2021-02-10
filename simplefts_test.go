// Copied from https://github.com/akrylysov/simplefts/blob/master/index_test.go

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	idx := make(index)

	assert.Nil(t, idx.search("foo"))
	assert.Nil(t, idx.search("donut"))

	idx.add([]Document{{ID: 1, Text: "A donut on a glass plate. Only the donuts."}})
	assert.Nil(t, idx.search("a"))
	assert.Equal(t, idx.search("donut"), []int{1})
	assert.Equal(t, idx.search("DoNuts"), []int{1})
	assert.Equal(t, idx.search("glass"), []int{1})

	idx.add([]Document{{ID: 2, Text: "donut is a donut"}})
	assert.Nil(t, idx.search("a"))
	assert.Equal(t, idx.search("donut"), []int{1, 2})
	assert.Equal(t, idx.search("DoNuts"), []int{1, 2})
	assert.Equal(t, idx.search("glass"), []int{1})
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
			assert.EqualValues(st, tc.tokens, tokenize(tc.text))
		})
	}
}

func TestLowercaseFilter(t *testing.T) {
	var (
		in  = []string{"Cat", "DOG", "fish"}
		out = []string{"cat", "dog", "fish"}
	)
	assert.Equal(t, out, lowercaseFilter(in))
}

func TestStopwordFilter(t *testing.T) {
	var (
		in  = []string{"i", "am", "the", "cat"}
		out = []string{"am", "cat"}
	)
	assert.Equal(t, out, stopwordFilter(in))
}

func TestStemmerFilter(t *testing.T) {
	var (
		in  = []string{"cat", "cats", "fish", "fishing", "fished", "airline"}
		out = []string{"cat", "cat", "fish", "fish", "fish", "airlin"}
	)
	assert.Equal(t, out, stemmerFilter(in))
}
