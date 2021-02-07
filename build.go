package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/BurntSushi/toml"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Config struct {
	Input  InputInfo
	Output OutputInfo
}

type InputInfo struct {
	BaseDir string `toml:"base_directory"`
	Files   []map[string]string
}

type OutputInfo struct {
	Dirname string
}

type Document struct {
	DocumentID int
	Path       string
	URL        string
	Title      string
}

type Index struct {
	Documents map[int]Document         `json:"documents"`
	Words     map[string]map[int][]int `json:"words"`
}

func main() {
	configPath := os.Args[1]
	configToml, err := ioutil.ReadFile(configPath)
	check(err)

	config := Config{}
	err = toml.Unmarshal(configToml, &config)
	check(err)

	everythingIndex := Index{
		Documents: map[int]Document{},
		Words:     map[string]map[int][]int{},
	}

	for documentID, file := range config.Input.Files {
		path, missingPath := file["path"]
		if missingPath == false {
			log.Fatal("Missing file path")
		}

		url, missingURL := file["url"]
		if missingURL == false {
			log.Fatal("Missing file url")
		}

		title, missingTitle := file["title"]
		if missingTitle == false {
			log.Fatal("Missing file title")
		}

		f, err := os.Open(path)
		check(err)

		r := bufio.NewReader(f)

		position := 0
		startPosition := 0
		word := make([]rune, 0)

		addWord := func(index Index, word string, path string, url string, title string, startPosition int) {
			index.Documents[documentID] = Document{DocumentID: documentID, Path: path, URL: url, Title: title}
			wordLower := strings.ToLower(word)
			if _, ok := index.Words[wordLower]; ok {
				if _, ok := index.Words[wordLower][documentID]; ok {
					index.Words[wordLower][documentID] = append(index.Words[wordLower][documentID], startPosition)
				} else {
					index.Words[wordLower][documentID] = []int{startPosition}
				}
			} else {
				index.Words[wordLower] = map[int][]int{}
			}
		}

		for {
			if character, _, err := r.ReadRune(); err != nil {
				if err == io.EOF {
					addWord(everythingIndex, string(word), path, url, title, startPosition)
					break
				} else {
					log.Fatal(err)
				}
			} else {
				if !unicode.IsLetter(character) {
					// Handle multiple spaces
					if len(word) == 0 {
						continue
					}

					addWord(everythingIndex, string(word), path, url, title, startPosition)
					startPosition = position + 1
					word = make([]rune, 0)
				} else {
					word = append(word, character)
				}
				position++
			}
		}
		_ = startPosition
	}

	splitIndexes := map[string]Index{}
	characters := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	for _, firstCharacter := range characters {
		for _, secondCharacter := range characters {
			ID := strings.Join([]string{firstCharacter, secondCharacter}, "")
			splitIndexes[ID] = Index{
				Documents: map[int]Document{},
				Words:     map[string]map[int][]int{},
			}
		}
	}

	for word, documents := range everythingIndex.Words {
		ID := shortHash(word, 2)
		splitIndexes[ID].Words[word] = documents
		for documentID := range documents {
			splitIndexes[ID].Documents[documentID] = everythingIndex.Documents[documentID]
		}
	}

	encode := func(dirname string, filename string, index Index) {
		if _, err := os.Stat(dirname); os.IsNotExist(err) {
			os.Mkdir(dirname, os.ModeDir)
		}

		indexPath := filepath.Join(dirname, filename)
		fmt.Println(indexPath)
		file, err := os.Create(indexPath)
		check(err)
		defer file.Close()
		encoder := gob.NewEncoder(file)
		encoder.Encode(index)
	}

	err = os.RemoveAll(config.Output.Dirname)
	check(err)

	for indexName, index := range splitIndexes {
		encode(config.Output.Dirname, indexName+".ltd", index)
	}
}

func shortHash(s string, size int) string {
	h := sha1.New()
	h.Write([]byte(s))
	return string(hex.EncodeToString(h.Sum(nil))[:size])
}
