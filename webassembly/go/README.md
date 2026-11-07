# A ÷ B calculator — Go → WebAssembly

This project compiles Go code to WebAssembly with the **official Go compiler**
and uses it from a web page.

## Run it

You only need Docker.

```sh
docker compose up --build
```

Open **http://localhost:8080**, enter two numbers and press **Divide**.
Try `B = 0` to see the error that comes from Go. Stop with `Ctrl+C`, or with
`docker compose down` if you started it with `-d`.

## See the built files on your disk

`main.wasm` and `wasm_exec.js` are created inside Docker. To get a copy of them in
the `web/` folder too:

```sh
docker build --target export --output web .
```

This runs the build and writes the two files from the Dockerfile's `export` stage
into `web/`. Run it again after you change `main.go`.

- `main.wasm` is binary. Its first four bytes are `\0asm`, the WebAssembly "magic
  number" (`head -c 4 web/main.wasm | od -c` shows them).
- `wasm_exec.js` is plain JavaScript you can open and read.

The copies are only for looking at. `.gitignore` keeps them out of git, and
`.dockerignore` keeps them out of the build, so an old copy in `web/` can never
replace the freshly built file.

## The files, in reading order

```
go/
├── src/
│   ├── main.go       1. the Go code: the division + the bridge to JavaScript
│   └── go.mod           (the usual Go module file)
├── Dockerfile        2. how main.wasm is built and served
├── web/
│   ├── index.html    3. the page
│   ├── app.js        4. loads main.wasm and calls the Go function
│   ├── style.css        (just looks)
│   ├── main.wasm        (built copy, only after the export command above)
│   └── wasm_exec.js     (built copy, only after the export command above)
├── compose.yaml         (starts the container on port 8080)
├── .dockerignore        (which files Docker may use for the build)
└── .gitignore           (keeps the built copies out of git)
```

### 1. `src/main.go` — the Go side

The whole Go program is one file with three parts:

- **`Divide(a, b)`**: the actual calculation. Plain Go, nothing
  WebAssembly-specific. It returns an error when `b` is zero.
- **`divide(this, args)`**: the bridge. JavaScript calls it with two numbers,
  which arrive as `js.Value`. `.Float()` turns them into Go numbers. It calls
  `Divide` and returns a Go `map`, which JavaScript sees as an object:
  `{ result: 3 }` or `{ error: "cannot divide by zero" }`.
- **`main()`**: `js.Global().Set("divide", js.FuncOf(divide))` puts the
  function on JavaScript's global `window` object. `select {}` then keeps the Go
  program running. If `main()` returned, Go would stop and `divide` could no
  longer be called.

The first line, `//go:build js && wasm`, says the file is meant only for the
browser target, because the `syscall/js` package exists only there.

### 2. `Dockerfile` — building and serving

Three stages:

1. **build** (image `golang`):
   - `GOOS=js GOARCH=wasm go build -o main.wasm` compiles for the browser instead
     of for Linux.
   - It copies **`wasm_exec.js`** from the Go installation. This is Go's JavaScript
     glue, and it must come from the same Go version that built `main.wasm`. That's
     why it's copied at build time and not committed to the repo.
2. **export** (optional): an empty image holding only `main.wasm` and
   `wasm_exec.js`. It's used only by the `--target export --output web` command
   above.
3. **serve** (image `nginx`): a small web server that hosts `main.wasm`,
   `wasm_exec.js` and the files from `web/`. It's the last stage, so it's what
   `docker compose up --build` builds and runs.

### 3. `web/index.html` — the page

A form with inputs A and B and a **Divide** button. The button stays disabled until
WebAssembly has loaded. At the bottom, two scripts load **in this order**:

1. `wasm_exec.js`, which defines the `Go` class.
2. `app.js`, our code, which uses it.

### 4. `web/app.js` — calling Go from JavaScript

1. `new Go()` creates the helper from `wasm_exec.js`.
2. `WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)` downloads
   and compiles the module. `go.importObject` contains the browser functions the Go
   program needs.
3. `go.run(instance)` starts Go's `main()`, which creates `divide`.
4. On submit, the inputs are checked (`,` is accepted as a decimal separator). Then
   `divide(a, b)` is called like any JavaScript function, and its `result` or
   `error` is shown.

Checking that the user typed numbers is the page's job. The rule "no dividing by
zero" is Go's job.

## What happens when you open the page

```
browser loads index.html
  → wasm_exec.js       defines class Go
  → app.js             fetches main.wasm, go.run() starts Go
       → Go main()     sets window.divide, then waits
  → button enabled     "WebAssembly loaded. Ready."
user clicks Divide
  → app.js             divide(7.5, 2.5)
       → Go divide()   → Divide() → returns { result: 3 }
  → page shows         7.5 ÷ 2.5 = 3
```

## Good to know

- **A web server is required.** Opening `index.html` directly from disk
  (`file://...`) doesn't work, because browsers block `fetch()` of local files. The
  server must also send `main.wasm` with the content type `application/wasm`.
  nginx does that by default.
- **Decimal numbers are floats.** Go's `float64` and JavaScript numbers are the same
  64-bit floating point type, so `1 ÷ 3 = 0.3333333333333333`, exactly as JavaScript
  would show it.
- **The file is big (~2 MB).** `main.wasm` contains the full Go runtime (garbage
  collector, goroutine scheduler) and the parts of the standard library it uses.
  Compare with the [`tinygo/`](../tinygo/README.md) version of the same code (~59 KB).
