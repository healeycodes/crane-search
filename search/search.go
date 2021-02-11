package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"syscall/js"
	"time"

	lib "github.com/healeycodes/crane-search"
)

func main() {
	// Define the function "Crane" in the JavaScript scope
	js.Global().Set("Crane", Crane())
	// Prevent the function from returning, which is required in a wasm module
	select {}
}

// Crane fetches an external resource by making a HTTP request from Go
func Crane() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Get the URL as argument
		// args[0] is a js.Value, so we need to get a string out of it
		indexURL := args[0].String()
		search := args[1].String()

		fmt.Println(indexURL)

		// Handler for the Promise
		// We need to return a Promise because HTTP requests are blocking in Go
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			reject := args[1]

			// Run this code asynchronously
			go func() {
				// Make the HTTP request
				res, err := http.DefaultClient.Get(indexURL)
				if err != nil {
					// Handle errors: reject the Promise if we have an error
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New(err.Error())
					reject.Invoke(errorObject)
					return
				}
				defer res.Body.Close()

				// Read the response body
				data, err := ioutil.ReadAll(res.Body)
				if err != nil {
					// Handle errors here too
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New(err.Error())
					reject.Invoke(errorObject)
					return
				}

				buf := bytes.NewBuffer(data)
				dec := gob.NewDecoder(buf)

				store := lib.Store{}
				if err := dec.Decode(&store); err != nil {
					log.Fatal(err)
				}

				start := time.Now()
				matchedIDs := store.Index.Search(search)
				log.Printf("Search found %d documents in %v", len(matchedIDs), time.Since(start))

				results := []lib.Result{}
				for _, id := range matchedIDs {
					results = append(results, store.Results[id])
					log.Printf("%s\n", store.Results[id].Title)
				}

				result, err := json.Marshal(results)
				if err != nil {
					fmt.Println(err)
					return
				}

				// "data" is a byte slice, so we need to convert it to a JS Uint8Array object
				arrayConstructor := js.Global().Get("Uint8Array")
				dataJS := arrayConstructor.New(len(result))
				js.CopyBytesToJS(dataJS, result)

				// Create a Response object and pass the data
				responseConstructor := js.Global().Get("Response")
				response := responseConstructor.New(dataJS)

				// Resolve the Promise
				resolve.Invoke(response)
			}()

			// The handler of a Promise doesn't return any value
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
