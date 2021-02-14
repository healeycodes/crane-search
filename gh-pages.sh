./build-index.sh federalist.toml
./build-browser.sh
mkdir -p docs
cp -a ./dist/. ./docs/
