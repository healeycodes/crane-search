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
      } else {
        // We're probably on Safari/an older browser
        // so use WebAssembly.instantiate instead (slower)
        const go = new Go();
        fetch(wasmPath)
          .then((response) => response.arrayBuffer())
          .then((bytes) => WebAssembly.instantiate(bytes, go.importObject))
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
