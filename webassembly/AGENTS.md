# AGENTS.md

Guidance for AI coding agents working in this repository.
This file is the project's memory: the rules, every decision made while designing it (with the reasons and the rejected alternatives), what happened during implementation, and the original plan.

Sections:
1. Project rules
2. Decision log
3. Background established during design
4. Implementation record
5. Browser check script
6. Original implementation plan (as approved, unchanged)

## 1. Project rules
- This repo is a **tutorial** for writing and using WebAssembly from Go. Code is read by humans first: keep it short, plain and well commented. No clever tricks, no extra features, no unnecessary file splits.
- **Everything runs in Docker** (compiling, serving, any tooling, even browser checks). Never install anything on the host machine.
- `go/` and `tinygo/` are two independent, self-contained projects. They are intentionally duplicated — do not introduce shared files between them.
- Docker image versions are pinned in each `Dockerfile`. If you change one, rebuild and update the `.wasm` sizes quoted in the READMEs (root, `go/`, `tinygo/`).
- `README.md` files are for human readers (the tutorial). Agent-facing material (rules, decisions, plans) belongs here, not in the READMEs.
- Before changing anything listed in the decision log below, check the reason recorded there. Don't reverse a decision without the user asking.

## 2. Decision log

Decisions in the order they were made with the user (design conversation, 2026-11-05/06).

1. **Purpose and constraints** (the user's original request): a simple web page (HTML, JavaScript, basic CSS) used as a **tutorial** on writing and using WebAssembly. The calculation runs in WebAssembly written in **Go**. All compilation and development runs in **Docker**; nothing installed on the machine.

2. **Operation: A ÷ B, not A + B.** The first request was a sum. The user changed it to division, so that the Go side can also return an **error**: a plain, simple string message when B is zero (`"cannot divide by zero"`). `0 ÷ 0` is also that error.

3. **Two compilers, two separate projects.** The first plan used only the standard Go compiler, with TinyGo mentioned in the README. After TinyGo was explained (see section 3), the user chose to do **both**: `go/` and `tinygo/` as separate folders and separate, self-contained projects. Both must still be good tutorials: simple and very clear to read.
   - Rejected: one project only (either compiler); one page that switches between both builds (more code, less clear).

4. **TinyGo uses `syscall/js`, the same bridge as standard Go**, chosen by the user.
   - Effect: the Go code is identical in both projects (apart from one comment). The lesson is "same code, different compiler, ~34× smaller `.wasm`".
   - Rejected: TinyGo `//export divide` as a real wasm export. Wasm exports can only return numbers, so returning the error string would need extra exports (`errorPtr()`, `errorLen()`) and JavaScript reading bytes from `instance.exports.memory.buffer` with `TextDecoder`. That's about 15 extra lines of trickier code; it teaches raw WebAssembly but hurts simplicity.

5. **One Go file per project (`src/main.go`), no unit tests.** The first plan had three Go files: `divide.go` (pure logic), `divide_test.go`, and `main.go` (the bridge). The split existed only because `syscall/js` compiles only for `GOOS=js GOARCH=wasm`, so tests run natively in Docker needed the logic in a separate file. The user questioned whether three files were necessary for such a simple app; they were OK with it only if it kept the tutorial easy to read.
   - Decision: one `main.go` read top to bottom: `Divide()` (plain Go) → `divide()` (bridge) → `main()`. Plus the two-line `go.mod`.
   - Reason: the unit tests would teach Go testing, not WebAssembly. Correctness is checked in a real browser instead (section 5).
   - Don't re-add test files or split `main.go` unless the user asks.

6. **Division of responsibility.** JavaScript (`app.js`) checks the input: empty or not a number shows "Please enter two numbers." without calling wasm, and `,` is accepted as a decimal separator. Go owns the division rule and its error string. Go returns a JS object, `{ result }` or `{ error }`.

7. **Decimal numbers are `float64`** (IEEE-754, the same as JavaScript numbers). `1 ÷ 3 = 0.3333333333333333` is expected and explained in the READMEs. No decimal or big-number library.

8. **Where the plan and rules live.** The user asked for the plan to be written to "readme.md or agent.md, whichever is the correct one", and committed before implementation.
   - Decision: `AGENTS.md` (for AI tools) plus a one-line `CLAUDE.md` (`@AGENTS.md`) so Claude Code loads it automatically. `README.md` stays the human tutorial.
   - The user later asked why both files exist, considered trimming `AGENTS.md`, then asked for a solution where **nothing from the design conversation is lost**. That's why this file has the full decision log, background and record, and keeps the original plan unchanged.

9. **Commits only when the user asks.** The user asked for: (a) a commit of the plan before implementing, and (b) a commit after implementation. No remote is configured; nothing has been pushed.

10. **`Divide()` and `divide()` keep their names.**
    - The user asked why the two names differ only in the first letter. The explanation given:
      - `Divide` is the plain-Go calculation; `divide` is the bridge whose signature `js.FuncOf` requires.
      - Go's uppercase = exported rule has no practical effect in `package main`.
      - JavaScript sees the function only because of `js.Global().Set("divide", ...)`; the JS name is just that string, unrelated to the Go name or its case.
    - Proposed and **declined by the user**: renaming to `divide` (logic) + `divideFromJS` (bridge). Don't rename unless the user asks.

11. **Built files can be exported to `web/` on disk.** The user wanted `main.wasm` and `wasm_exec.js` on disk, not only inside Docker, to see the Go vs TinyGo differences instantly.
    - Three options were presented:
      1. A Dockerfile `export` stage plus `docker build --target export --output web .`
      2. `docker compose cp` from the running container.
      3. A bind-mount setup: a builder service compiles into `./web` and nginx serves `./web`.
    - The user chose **option 1**. It keeps the existing setup and is one command, and it teaches Docker as a build tool.
    - Rejected: option 2 (manual, needs a running container, easy to forget) and option 3 (rewrites the Docker part of both projects; more moving parts for readers).
    - Trap handled: the build stage runs `COPY web/ /site/` *after* compiling, so a stale `web/main.wasm` on disk would overwrite the fresh one. `.dockerignore` therefore excludes `web/main.wasm` and `web/wasm_exec.js`.
    - The copies are **not committed**: each project's `.gitignore` lists them. They're generated, `main.wasm` is 2 MB for Go, and `wasm_exec.js` must match the compiler version.
    - `serve` must stay the **last** Dockerfile stage, because `docker compose` builds the last stage by default (`compose.yaml` sets no `target`).
    - Mentioned but not implemented: reading `main.wasm` as text with `wasm2wat` (WABT) in a container.

## 3. Background established during design

The user asked "what is TinyGo?" and then whether TinyGo is "the proper way" and standard Go "non-standard". These are the answers given, and they shaped decisions 3 and 4:

- **TinyGo** is a second, independent Go compiler built on LLVM, originally for microcontrollers and WebAssembly. Compared with standard Go, its wasm output is much smaller (here 58,948 vs 2,027,861 bytes). It supports almost the whole language, but some standard library packages are partial (mostly around `reflect`, for example `encoding/json`). It supports new Go versions a little later. Its GC is simpler and can be slower when a program allocates heavily. It ships its **own** `wasm_exec.js`, which is not interchangeable with Go's.
- **Neither compiler is non-standard.** Both emit valid, standard WebAssembly modules. The WebAssembly standard defines the module format, imports and exports, and numeric types only (`i32`, `i64`, `f32`, `f64`). It does **not** define how strings, objects, GC or errors cross into JavaScript, so every language ships its own glue: Go and TinyGo have `wasm_exec.js` + `syscall/js`, Rust has `wasm-bindgen`, C/C++ has Emscripten's JS.
- **Standard Go (`GOOS=js`)** runs the full Go runtime inside the module and talks to JS through imports provided by `wasm_exec.js`. Functions become callable because Go sets them on the JS global object, not because they appear in `instance.exports`. This is a Go convention on top of standard wasm, and it's the official, supported route.
- **TinyGo `//export`** creates real wasm exports with numbers-only signatures. That's the thinnest layer over the raw standard and teaches the core mechanism more directly. It's a teaching advantage, not "more correct".
- Go 1.24+ also has `//go:wasmexport`, mainly aimed at the WASI target (`GOOS=wasip1`). Not used here.

## 4. Implementation record

Implemented in commit `5b6755d` (plan committed earlier in `c0d9c19`).

- **Pinned images** (latest stable at implementation time, 2026-11-06):
  - `golang:1.27.3-alpine`
  - `tinygo/tinygo:0.44.0` (bundles go1.27.3 and LLVM 22.1.4)
  - `nginx:1.30.6-alpine` (1.30 is nginx's stable branch)
  - `go.mod` says `go 1.25`; this is a minimum, and nothing newer is needed.
- **`wasm_exec.js` locations:** Go: `$(go env GOROOT)/lib/wasm/wasm_exec.js`. TinyGo: `$(tinygo env TINYGOROOT)/targets/wasm_exec.js` (`/usr/local/tinygo`). Both are copied at build time and never committed, because each must match the compiler version that built `main.wasm`.
- **TinyGo image runs as user `tinygo` (uid 1000)**, so `tinygo/Dockerfile` has `USER root` in the build stage to create `/src` and `/site`. This deviates from the plan only in making it explicit.
- **nginx's default `mime.types` already maps `wasm` → `application/wasm`**, so no custom nginx config.
- **Build flags:** Go `GOOS=js GOARCH=wasm go build`. TinyGo `tinygo build -target wasm -no-debug`.
- **Measured `main.wasm` sizes:** Go 2,027,861 bytes (~2 MB); TinyGo 58,948 bytes (~59 KB), about 34× smaller.
- **Ports:** Go 8080, TinyGo 8081, so both can run at once.
- **Export stage (decision 11)**, added after the initial implementation.
  - Verified: `docker build --target export --output web .` writes only `main.wasm` and `wasm_exec.js` into `web/`, leaving the other web files untouched. Sizes: Go 2,027,861 + 16,992 bytes; TinyGo 58,948 + 18,002 bytes. Both files start with the `\0asm` magic bytes.
  - Verified: with junk deliberately written to `web/main.wasm`, `docker compose up --build` still served the real files (2,027,861 / 58,948 bytes), and the browser check passed on both ports.
- **Verification performed:**
  - Both servers return `200` with `Content-Type: application/wasm` for `main.wasm`, and serve `index.html`, `app.js`, `wasm_exec.js` and `style.css`.
  - The browser check (section 5) passed all 8 cases on both ports, with no console errors.
  - A screenshot was reviewed: centered card, error text in red.
- **Deviations from the plan:**
  - Versions are newer than the plan's examples (it said to pin whatever exists at implementation time).
  - The browser check used Playwright in Docker, not claude-in-chrome.
  - The README wording "compiling, testing and serving" became "compiling and serving", since there are no tests.
  - Rules about version pinning and the Playwright check were added to section 1.
- **Known, accepted edge cases (not raised with the user):**
  - JavaScript `Number()` accepts inputs like `Infinity` and `1e3`, so they're passed to Go as numbers.
  - Only the first `,` is replaced with `.`, so `1,000.5` is rejected as not a number.

## 5. Browser check script

Used to verify both projects (decision 5 replaced unit tests with this). To re-run it:

1. Start both projects (`docker compose up --build -d` in `go/` and `tinygo/`).
2. Save the script below as `check.py` in an **empty temporary folder** (not in this repo).
3. Run, from that folder:
   ```sh
   docker run --rm -v "$PWD":/out mcr.microsoft.com/playwright/python:v1.64.1-noble \
     sh -c 'pip install -q --break-system-packages playwright==1.64.1 && python3 /out/check.py'
   ```
   The image contains the browsers but not the `playwright` Python package. The pip version must match the image tag. Screenshots are written to the same folder.

```python
from playwright.sync_api import sync_playwright

CASES = [
    ("7.5", "2.5", "7.5 ÷ 2.5 = 3", False),
    ("1", "3", "1 ÷ 3 = 0.3333333333333333", False),
    ("-1", "4", "-1 ÷ 4 = -0.25", False),
    ("2,5", "2", "2.5 ÷ 2 = 1.25", False),
    ("5", "0", "cannot divide by zero", True),
    ("0", "0", "cannot divide by zero", True),
    ("5", "", "Please enter two numbers.", True),
    ("abc", "2", "Please enter two numbers.", True),
]

failed = 0
with sync_playwright() as p:
    browser = p.chromium.launch()
    for port in (8080, 8081):
        page = browser.new_page()
        errors = []
        page.on("console", lambda m: m.type == "error" and errors.append(m.text))
        page.on("pageerror", lambda e: errors.append(str(e)))
        page.goto(f"http://host.docker.internal:{port}/")
        page.wait_for_function("!document.getElementById('divide').disabled", timeout=10000)
        print(f"== port {port}: status = {page.text_content('#status')!r}")
        for a, b, expected, is_error in CASES:
            page.fill("#a", a)
            page.fill("#b", b)
            page.click("#divide")
            text = page.text_content("#output")
            err = "error" in (page.get_attribute("#output", "class") or "")
            ok = text == expected and err == is_error
            failed += not ok
            print(f"  {'OK  ' if ok else 'FAIL'} A={a!r:7} B={b!r:5} -> {text!r} (error style: {err})")
        print(f"  console errors: {errors or 'none'}")
        failed += len(errors)
        page.screenshot(path=f"/out/port{port}.png")
    browser.close()
print("ALL PASSED" if failed == 0 else f"{failed} FAILURE(S)")
```

## 6. Original implementation plan (as approved, unchanged)

Kept word for word as the user approved it. Where it differs from what was built, section 4 is correct.

### Implementation plan: WebAssembly division calculator tutorial — Go and TinyGo, Docker-only

#### Context
Goal: a tutorial showing how to write and use WebAssembly from Go. A plain HTML/CSS/JS page divides two decimal numbers, `A / B`. The division runs in WebAssembly compiled from Go. If B is zero, the Go code returns a plain string error message that the page shows.

Two **separate, self-contained projects** compile the same idea:
- `go/` — official Go compiler (`GOOS=js GOARCH=wasm`)
- `tinygo/` — TinyGo compiler

Both use the same bridge, `syscall/js` (user's choice). The Go code is nearly identical, so the takeaway is: same code, different compiler, very different `.wasm` size.

Constraints: everything (compiling, testing, serving) runs in Docker; nothing installed on the host (Docker 29.8 is available). The code is meant to be read by people, so it should be short, plain, well commented, with no clever tricks or extra features.

#### Layout
```
.
├── README.md                 # overview: what WebAssembly is, the two projects, comparison, how to run each
├── go/
│   ├── README.md             # step-by-step tutorial for this project
│   ├── Dockerfile            # stages: build → serve
│   ├── compose.yaml          # http://localhost:8080
│   ├── .dockerignore
│   ├── src/
│   │   ├── go.mod
│   │   └── main.go           # the whole Go side: Divide + JS bridge (~40 lines)
│   └── web/
│       ├── index.html
│       ├── style.css
│       └── app.js
└── tinygo/                   # same structure; compose on http://localhost:8081
    ├── README.md, Dockerfile, compose.yaml, .dockerignore
    ├── src/  (go.mod, main.go)
    └── web/  (index.html, style.css, app.js)
```
The two projects are intentionally duplicated (no shared files), so each folder can be read, copied and run alone. Different ports let both run at the same time for comparison.

#### Code (same in both projects unless noted)

**`src/main.go`** — the only Go file, read top to bottom (decision: one file, no unit tests; the tests would teach Go testing, not WebAssembly):
```go
//go:build js && wasm

package main

import (
	"errors"
	"syscall/js"
)

// Divide is the actual calculation: plain Go, nothing WebAssembly-specific.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

// divide is the bridge: JavaScript calls it, it calls Divide.
func divide(this js.Value, args []js.Value) any {
	result, err := Divide(args[0].Float(), args[1].Float())
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	return map[string]any{"result": result}
}

func main() {
	js.Global().Set("divide", js.FuncOf(divide))
	select {} // keep Go running so JavaScript can keep calling divide
}
```
Final version gets slightly richer explanatory comments (what `js.FuncOf` does, why `map` becomes a JS object, why `select {}`), but keeps this shape. Same file in both projects.

**`web/index.html`** — inputs A and B (`type="text" inputmode="decimal"`), a "÷" / "Divide" button (disabled until wasm loads), a result line, and a status line ("Loading WebAssembly…"). Loads `wasm_exec.js`, then `app.js`.

**`web/app.js`** — about 30 lines, in this order:
1. `const go = new Go();`
2. `WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)` → `go.run(instance)` → enable the button and show "Ready". On failure, show the error.
3. On click: read inputs, treat `,` as `.`, convert with `Number()`. If an input is empty or not a number, show a message in JS without calling wasm (comment: input validation is the page's job; the division rule is Go's job).
4. `const response = divide(a, b);` → show `response.error` (red) or `response.result`.

**`web/style.css`** — small centered card, readable font, light/dark using `prefers-color-scheme`, red error text.

#### Docker

**`go/Dockerfile`** (`golang:1.25-alpine`, `nginx:1.27-alpine`; pin exact tags that exist at implementation time):
1. `build`: `GOOS=js GOARCH=wasm go build -o /site/main.wasm .`. Copy `wasm_exec.js` from `$(go env GOROOT)/lib/wasm/` (Go ≥1.24 location; `misc/wasm` on older versions), and copy `web/*` → `/site`.
2. `serve` (final stage): nginx serving `/site`. nginx's built-in `mime.types` maps `.wasm` → `application/wasm`, which `instantiateStreaming` requires. Verify this and add an explicit mapping only if it's missing.

**`tinygo/Dockerfile`** (`tinygo/tinygo:<pinned version>`, `nginx:1.27-alpine`):
1. `build`: `tinygo build -o /site/main.wasm -target wasm -no-debug .`. Copy TinyGo's **own** `wasm_exec.js` from `$(tinygo env TINYGOROOT)/targets/wasm_exec.js` (it's not interchangeable with Go's, and the README says so). Run the stage as root if the image's default user can't write to `/site`.
2. `serve`: same as Go.

Each Dockerfile is commented line by line, since it is part of the tutorial.

**`compose.yaml`** — one `web` service: `build: .`, ports `8080:80` (Go) or `8081:80` (TinyGo).

#### READMEs
- **Root `README.md`**: what WebAssembly is (module, imports and exports, numbers-only ABI, why each language ships JS glue); the two projects; a comparison table (compiler, image, `.wasm` size filled in with real measured numbers, glue file); quick start for each.
- **`go/README.md` and `tinygo/README.md`**: prerequisites (Docker only) · run · files explained in reading order (`main.go` → `Dockerfile` → `index.html` → `app.js`) · what happens at runtime from page load to result · why a web server is needed (no `file://`, `.wasm` MIME type) · floating-point note (`1 / 3 = 0.3333333333333333`, IEEE-754 float64, same as JavaScript). The TinyGo README also covers what TinyGo is, its separate `wasm_exec.js`, and its size advantage and limits.

#### Commands (all through Docker)
- `cd go && docker compose up --build` → http://localhost:8080
- `cd tinygo && docker compose up --build` → http://localhost:8081

#### Verification
1. `docker compose up --build -d` in both folders (the build itself proves both compilers accept `main.go`). Then:
   - `curl -sI localhost:8080/main.wasm` and `curl -sI localhost:8081/main.wasm` return `200` with `Content-Type: application/wasm`.
   - `index.html`, `app.js` and `wasm_exec.js` are all served.
2. Record both `.wasm` sizes (from `docker run --rm <image> ls -l /usr/share/nginx/html`) for the README comparison table.
3. In a browser (claude-in-chrome if available, otherwise ask the user), on both ports:
   - `7.5 / 2.5` → `3`
   - `1 / 3` → `0.3333333333333333`
   - `-1 / 4` → `-0.25`
   - `5 / 0` → "cannot divide by zero"
   - empty B → validation message
   - no console errors
4. `docker compose down` in both folders.
5. Commit only if the user asks.
