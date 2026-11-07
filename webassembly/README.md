# WebAssembly with Go — a tutorial

A tiny web calculator that divides two numbers, **A ÷ B**.
The page is plain HTML, CSS and JavaScript, but the division itself runs in
**WebAssembly** compiled from **Go**. If B is zero, the Go code returns an error
message, which the page shows.

The same calculator is built twice, with two different compilers:

| Folder | Compiler | Port | Size of `main.wasm` |
|---|---|---|---|
| [`go/`](go/README.md) | the official Go compiler | http://localhost:8080 | ~2 MB (2,027,861 bytes) |
| [`tinygo/`](tinygo/README.md) | TinyGo | http://localhost:8081 | ~59 KB (58,948 bytes) |

The Go code is the same in both. Only the compiler changes, and the TinyGo
file is about **34× smaller**.

Each folder is a complete project on its own. You can read, copy or run one without
the other. Start with [`go/`](go/README.md), then look at what changes in
[`tinygo/`](tinygo/README.md).

## You only need Docker

Compiling and serving both run inside Docker containers. You don't install
Go, TinyGo or a web server on your computer.

```sh
cd go                       # or: cd tinygo
docker compose up --build   # build and start
# open http://localhost:8080 (go) or http://localhost:8081 (tinygo)
docker compose down         # stop when you are done

# optional: copy the built main.wasm and wasm_exec.js into web/ to look at them
docker build --target export --output web .
```

With both projects exported, you can compare them directly:
`ls -l go/web/main.wasm tinygo/web/main.wasm` and
`diff go/web/wasm_exec.js tinygo/web/wasm_exec.js`.

## What is WebAssembly?

WebAssembly (**wasm**) is a compact binary format that every modern browser can run,
next to JavaScript. You don't usually write it by hand. You write code in a language
such as Go, Rust or C, and a compiler turns it into a `.wasm` file.

A few ideas are enough to follow this tutorial:

- **Module**: a `.wasm` file. The browser downloads it and runs it in a sandbox, as
  safely as JavaScript.
- **Exports and imports**: a module can offer functions to JavaScript (exports) and
  can use functions that JavaScript gives it (imports). That's the only way the
  two sides talk to each other.
- **Only numbers**: at the lowest level, WebAssembly functions only take and return
  numbers (integers and floats). Strings, objects and errors need extra code to
  cross between JavaScript and wasm.
- **Glue code**: because of that, every language ships a small JavaScript helper
  that handles the details. For Go and TinyGo this helper is `wasm_exec.js`, and
  Go's `syscall/js` package uses it to read JavaScript values and return a JavaScript
  object such as `{ result: 2.5 }` or `{ error: "cannot divide by zero" }`.

## How the pieces fit together

```
 index.html ──loads──▶ wasm_exec.js   (glue from Go/TinyGo)
     │
     └──loads──▶ app.js ──fetches──▶ main.wasm   (your Go code, compiled)
                   │                     │
                   │   go.run() starts Go; Go's main() puts divide() on window
                   │                     │
                   └──── divide(a, b) ───┘──▶ { result } or { error }
```
