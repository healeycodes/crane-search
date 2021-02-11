GOOS=js GOARCH=wasm go build -o dist/search.wasm search/search.go
cp dist/search.wasm demo/search.wasm
