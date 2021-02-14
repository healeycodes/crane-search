package main

import (
	"encoding/gob"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	search "github.com/healeycodes/crane-search"
)

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

func main() {
	build()
}

// Builds the index and results into a store.
// Compresses it, encodes it, and writes it to disk
func build() {
	if len(os.Args) < 2 {
		log.Fatalln("Missing index argument")
		return
	}

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
	index := make(search.Index)
	index.Add(documents)
	log.Printf("Indexed %d documents in %v", len(documents), time.Since(start))

	// Strip text from documents to create results (so we don't send the full text)
	results := []search.Result{}
	for _, document := range documents {
		result := search.Result{
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
	store := search.Store{
		Index:   index,
		Results: results,
	}
	encoder.Encode(store)
}

// Given a config, load the items and their metadata.
func loadDocuments(config config) ([]search.Document, error) {
	docs := []search.Document{}
	for id, file := range config.Input.Files {
		path, missingPath := file["path"]
		url, missingURL := file["url"]
		title, missingTitle := file["title"]
		if !missingPath || !missingURL || !missingTitle {
			return nil, errors.New("Missing metadata")
		}

		data, err := ioutil.ReadFile(path)
		check(err)
		docs = append(docs, search.Document{
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
