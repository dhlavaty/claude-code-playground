# A ÷ B calculator — TinyGo → WebAssembly

This project is the same calculator as [`go/`](../go/README.md), but the Go code is
compiled with **TinyGo** instead of the official Go compiler.

If you haven't read the [`go/` tutorial](../go/README.md) yet, start there. This
page only explains what is different.

## Run it

You only need Docker.

```sh
docker compose up --build
```

Open **http://localhost:8081**, enter two numbers and press **Divide**.
(It uses port 8081, so you can run it next to the `go/` version on 8080.)

## What is TinyGo?

TinyGo is a second, independent compiler for the Go language. It's built on LLVM
and was made for small targets: microcontrollers and WebAssembly. You write normal
Go code, but compile it with `tinygo build` instead of `go build`.

| | Go | TinyGo |
|---|---|---|
| Who makes it | the Go team (official) | a separate open-source project |
| Size of our `main.wasm` | ~2 MB (2,027,861 bytes) | **~59 KB (58,948 bytes)** |
| Go language support | everything | almost everything; some standard library packages are incomplete |
| New Go versions | supported immediately | supported a bit later |

TinyGo's files are smaller because its runtime is much simpler and it removes all
code the program doesn't use.

## What is different from the `go/` project

Only three things. Compare the files side by side to see them.

### 1. `src/main.go` — no change

The Go code is the same, apart from one comment. `syscall/js` works the same way in
TinyGo.

### 2. `Dockerfile` — a different compiler

```dockerfile
FROM tinygo/tinygo:0.44.0 AS build
...
RUN tinygo build -o /site/main.wasm -target wasm -no-debug .
```

- `-target wasm` means "WebAssembly for a web page", the same idea as
  `GOOS=js GOARCH=wasm` in standard Go.
- `-no-debug` leaves out debug information, which makes the file smaller.
- `USER root` is needed only because the TinyGo image runs as a normal user, and
  we want to create the `/src` and `/site` folders.

### 3. `wasm_exec.js` — TinyGo's own version

```dockerfile
RUN cp "$(tinygo env TINYGOROOT)/targets/wasm_exec.js" /site/
```

TinyGo ships its **own** `wasm_exec.js`. It has the same name and the same
`new Go()` / `go.run()` usage, but the inside is different. Using Go's
`wasm_exec.js` with a TinyGo `main.wasm`, or the other way round, doesn't work.
Always copy it from the compiler that built `main.wasm`.

`index.html` and `app.js` are otherwise the same as in the `go/` project.

## See the built files on your disk

Exactly as in the `go/` project, this writes `main.wasm` and `wasm_exec.js` into
`web/`:

```sh
docker build --target export --output web .
```

Run it in both projects, then compare them from the repository root:

```sh
ls -l go/web/main.wasm tinygo/web/main.wasm        # ~2 MB vs ~59 KB
diff go/web/wasm_exec.js tinygo/web/wasm_exec.js   # two different glue files
```

## When to choose which

- **Go**: the official toolchain. Every package and feature works, and it gets the
  newest Go version right away. The price is a bigger download.
- **TinyGo**: much smaller files, so the page loads faster. A good fit for small,
  focused pieces of logic like this one. Check that the packages you need are
  supported (see the TinyGo docs: https://tinygo.org/docs/reference/lang-support/stdlib/).
