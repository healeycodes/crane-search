// Uses code from Simple Full-Text Search engine
// Copied/Modified from https://github.com/akrylysov/simplefts
// LICENSE: https://github.com/akrylysov/simplefts/blob/master/LICENSE

package lib

import (
	"encoding/gob"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/BurntSushi/toml"
	snowballeng "github.com/kljensen/snowball/english"
)

// Document represents a text file
type Document struct {
	Title string
	URL   string
	Text  string
	ID    int
}

// Result represents a search result item
type Result struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	ID    int    `json:"id"`
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

type config struct {
	Input  inputInfo
	Output outputInfo
}

type inputInfo struct {
	BaseDir string `toml:"base_directory"`
	Files   []map[string]string
}

type outputInfo struct {
	Filename string
}

// Build builds the index, compresses it, encodes it, and writes it to disk
func Build() {
	configPath := os.Args[1]
	configToml, err := ioutil.ReadFile(configPath)
	check(err)

	config := config{}
	err = toml.Unmarshal(configToml, &config)
	check(err)

	start := time.Now()
	documents, err := loadDocuments(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d documents in %v", len(documents), time.Since(start))

	start = time.Now()
	index := make(Index)
	index.Add(documents)
	log.Printf("Indexed %d documents in %v", len(documents), time.Since(start))

	// Strip text from documents to create results
	results := []Result{}
	for _, document := range documents {
		result := Result{
			Title: document.Title,
			URL:   document.URL,
			ID:    document.ID,
		}
		results = append(results, result)
	}

	removeIfExists(config.Output.Filename)
	file, err := os.Create(config.Output.Filename)
	check(err)
	defer file.Close()

	encoder := gob.NewEncoder(file)
	store := Store{
		Index:   index,
		Results: results,
	}
	encoder.Encode(store)
}

func loadDocuments(config config) ([]Document, error) {
	docs := []Document{}
	for id, file := range config.Input.Files {
		path, missingPath := file["path"]
		url, missingURL := file["url"]
		title, missingTitle := file["title"]
		if !missingPath || !missingURL || !missingTitle {
			return nil, errors.New("Missing metadata")
		}

		data, err := ioutil.ReadFile(path)
		check(err)
		docs = append(docs, Document{
			Title: title,
			URL:   url,
			Text:  string(data),
			ID:    id,
		})
	}

	return docs, nil
}

func removeIfExists(path string) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return
		}
	}
	err := os.Remove(path)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
