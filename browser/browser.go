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
	js.Global().Set("_craneLoad", js.FuncOf(load))
	js.Global().Set("_craneQuery", js.FuncOf(query))
	select {}
}

func load(this js.Value, args []js.Value) interface{} {
	b := make([]byte, args[0].Get("length").Int())
	js.CopyBytesToGo(b, args[0])
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&store); err != nil {
		return js.ValueOf(err.Error())
	}
	return js.ValueOf(true)
}

func query(this js.Value, args []js.Value) interface{} {
	searchTerm := args[0].String()
	start := time.Now()
	matchedIDs := store.Index.Search(searchTerm)
	log.Printf("Search found %d documents in %v", len(matchedIDs), time.Since(start))

	var results []interface{}
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
