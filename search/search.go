package search

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"syscall/js"
)

// Main function: it sets up our Wasm application
func main() {
	// Define the function "LongTailedDuck" in the JavaScript scope
	js.Global().Set("LongTailedDuck", LongTailedDuck())
	// Prevent the function from returning, which is required in a wasm module
	select {}
}

// LongTailedDuck fetches an external resource by making a HTTP request from Go
func LongTailedDuck() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Get the URL as argument
		// args[0] is a js.Value, so we need to get a string out of it
		indexURL := args[0].String()
		word := args[1].String()

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

				// Decompress

				ungzipper, err := gzip.NewReader(res.Body)
				uncompressed := []byte{}
				_, err = ungzipper.Read(uncompressed)

				buf := bytes.NewBuffer(uncompressed)
				dec := gob.NewDecoder(buf)

				index := Index{}
				if err := dec.Decode(&index); err != nil {
					log.Fatal(err)
				}

				matchingDocuments := []Result{}
				documentInfo := index.Words[word]
				for documentID, positions := range documentInfo {
					matchingDocuments = append(matchingDocuments, Result{Document: index.Documents[documentID], Positions: positions})
				}

				result, err := json.Marshal(matchingDocuments)
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

func shortHash(s string, size int) string {
	h := sha1.New()
	h.Write([]byte(s))
	return string(hex.EncodeToString(h.Sum(nil))[:size])
}
