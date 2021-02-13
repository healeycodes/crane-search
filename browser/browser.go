package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"syscall/js"
	"time"

	search "github.com/healeycodes/crane-search"
)

var store search.Store

func main() {
	js.Global().Set("_craneLoad", js.FuncOf(Load))
	js.Global().Set("_craneSearch", js.FuncOf(Search))
	select {}
}

func Load(this js.Value, args []js.Value) interface{} {
	storeBytes := []byte{}
	for i := 0; i < args[0].Length(); i++ {
		storeBytes = append(storeBytes, byte(args[0].Index(i).Int()))
	}
	// js.CopyBytesToGo(storeBytes, args[0])
	buf := bytes.NewBuffer(storeBytes)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&store); err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func Search(this js.Value, args []js.Value) interface{} {
	searchTerm := args[0].String()
	start := time.Now()
	matchedIDs := store.Index.Search(searchTerm)
	log.Printf("Search found %d documents in %v", len(matchedIDs), time.Since(start))

	// results := map[string]interface{}{}
	results := make([]interface{}, len(matchedIDs))
	for _, id := range matchedIDs {
		results = append(results,
			map[string]interface{}{
				"title": store.Results[id].Title,
				"url":   store.Results[id].URL,
				"id":    store.Results[id].ID,
			})
	}

	return js.ValueOf(results)
}
