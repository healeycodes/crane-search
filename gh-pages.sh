./build-index.sh federalist.toml
./build-browser.sh
mkdir -p dist
mkdir -p docs
cp -a ./dist/. ./demo/
cp -a ./demo/. ./docs/
