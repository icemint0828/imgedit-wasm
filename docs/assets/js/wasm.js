const go = new Go(); // Defined in wasm_exec.js
const WASM_URL = 'main.wasm';
let wasm, mod;

// Polyfill without WebAssembly.instantiateStreaming
if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}
// Load the built Go program in main.wasm
WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then(function (obj) {
    wasm = obj.instance;
    mod = obj.module;

    go.run(wasm)

    // Enable Execute button
    let elements = document.getElementsByClassName("button")

    for( let i = 0 ; i < elements.length ; i ++ ) {
        elements[i].disabled = false
    }
    // Defined by main.wasm
    setValidValues()
})