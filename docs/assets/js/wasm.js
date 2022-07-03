const go = new Go(); // Defined in wasm_exec.js
const WASM_URL = 'main.wasm';
let wasm, mod;

// WebAssembly.instantiateStreamingがない場合のポリフィル
if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}
// main.wasmにビルドされたGoのプログラムを読み込む
WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then(function (obj) {
    wasm = obj.instance;
    mod = obj.module;

    go.run(wasm)

    // 実行ボタンを有効にする
    let elements = document.getElementsByClassName("button")

    for( let i = 0 ; i < elements.length ; i ++ ) {
        elements[i].disabled = false
    }
})