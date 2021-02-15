if (!WebAssembly.instantiateStreaming) {
  // polyfill
  // https://github.com/golang/go/blob/b2fcfc1a50fbd46556f7075f7f1fbf600b5c9e5d/misc/wasm/wasm_exec.html#L17-L22
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

class Crane {
  loadWasm(wasmPath, store) {
    return new Promise((resolve, reject) => {
      if (WebAssembly.instantiateStreaming !== undefined) {
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch(wasmPath), go.importObject)
          .then((result) => {
            go.run(result.instance);
            _craneLoad(store);
            resolve();
          })
          .catch((err) => reject(err));
      }
    });
  }
  async loadResource(storePath) {
    const body = await fetch(storePath).then((res) => res.arrayBuffer());
    return new Uint8Array(body);
  }
  query(searchTerm) {
    if (window._craneQuery === undefined) {
      console.warn("query: called before Crane has loaded");
      return [];
    }
    return _craneQuery(searchTerm);
  }
  async load() {
    try {
      const store = await this.loadResource(this.storePath);
      return await this.loadWasm(this.wasmPath, store);
    } catch (error) {
      console.error(`load: ${error}`);
    }
  }
  constructor(wasmPath, storePath) {
    this.wasmPath = wasmPath;
    this.storePath = storePath;
  }
}
