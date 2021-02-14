mkdir -p dist
GOOS=js GOARCH=wasm go build -o dist/crane.wasm browser/browser.go
