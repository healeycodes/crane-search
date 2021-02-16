# Crane 🐦

Crane is a technical demo is inspired by [Stork](https://github.com/jameslittle230/stork) and uses a near-identical configuration file setup. So it had to be named after a bird too.

I wrote it to help me understand how WebAssembly search tools work.

Crane is two programs. The first program scans a group of documents and builds an efficient index. 1MB of text and metadata is turned into a 25KB index (14KB gzipped). The second program is a Wasm module that is sent to the browser along with a little bit of JavaScript glue code and the index. The result is an instant search engine that helps users find web pages as they type.

![Crane instant search in action](https://github.com/healeycodes/crane-search/blob/main/docs/crane.gif)

The full text search engine is powered in part with code from Artem Krylysov's blog post [Let's build a Full-Text Search engine](https://artem.krylysov.com/blog/2020/07/28/lets-build-a-full-text-search-engine/).

## Use it

Describe your document files and their metadata.

```toml
[input]
files = [
    {
        path = "docs/essays/essay01.txt",
        url = "essays/essay01.txt",
        title = "Introduction"
    },
    # etc.
]

[output]
filename = "dist/federalist.crane"
```

Pass the configuration file to the build script. You'll want a fresh index whenever your documents change but you only need to build the Wasm module once ever.

```bash
./build-index.sh federalist.toml
./build-search.sh
```

Copy the files from `/dist` to your website. And away you go!

```javascript
const crane = new Crane("crane.wasm", "federalist.crane");
await crane.load()

const results = crane.query('some keywords');
console.log(results);
```

See the demo inside `/docs` for a basic UI.

<br>

## Build demo page

```bash
./gh-pages.sh
```
