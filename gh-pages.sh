./build-index.sh federalist.toml
./build-search.sh
mkdir -p dist
mkdir -p docs
cp -a ./dist/. ./demo/
cp -a ./demo/. ./docs/
