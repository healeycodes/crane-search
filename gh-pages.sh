./build-index.sh federalist.toml
./build-search.sh
cp -a ./dist/. ./demo/
cp -a ./demo/. ./docs/
