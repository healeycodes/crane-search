GOOS=js GOARCH=wasm go build -o dist/crane.wasm browser/browser.go
cp dist/crane.wasm demo/crane.wasm
