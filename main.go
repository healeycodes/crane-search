package main

import (
	"compress/gzip"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
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

type document struct {
	Title string
	URL   string
	Text  string
	ID    int
}

type Result struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	ID    int    `json:"id"`
}

func loadDocuments(config config) ([]document, error) {
	docs := []document{}
	for id, file := range config.Input.Files {
		path, missingPath := file["path"]
		url, missingURL := file["url"]
		title, missingTitle := file["title"]
		if !missingPath || !missingURL || !missingTitle {
			return nil, errors.New("Missing metadata")
		}

		data, err := ioutil.ReadFile(path)
		check(err)
		docs = append(docs, document{
			Title: title,
			URL:   url,
			Text:  string(data),
			ID:    id,
		})
	}

	return docs, nil
}

func main() {
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

	// start = time.Now()
	// matchedIDs := index.search("the States")
	// log.Printf("Search found %d documents in %v", len(matchedIDs), time.Since(start))

	// for _, id := range matchedIDs {
	// 	doc := docs[id]
	// 	log.Printf("%d\t%s\n", id, doc.Text)
	// }

	removeIfExists(config.Output.Filename)
	file, err := os.Create(config.Output.Filename)
	check(err)
	defer file.Close()

	gzipper, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	check(err)
	defer gzipper.Close()

	encoder := gob.NewEncoder(gzipper)
	encoder.Encode(results)
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
