# This Dockerfile has three stages:
#   1. "build"  - compiles the Go code into main.wasm using TinyGo
#   2. "export" - only main.wasm and wasm_exec.js, so you can copy them to your disk
#   3. "serve"  - a small web server that hosts the page and main.wasm
# The last stage (serve) is the default, so only it ends up in the final image.

# ---------------------------------------------------------------
# Stage 1: build
# ---------------------------------------------------------------
FROM tinygo/tinygo:0.44.0 AS build

# The TinyGo image runs as a normal user by default.
# Switch to root so we are allowed to create the /src and /site folders.
USER root

WORKDIR /src
COPY src/ .

# Compile for the browser:
#   -target wasm -> produce WebAssembly for a web page (same idea as GOOS=js GOARCH=wasm)
#   -no-debug    -> leave out debug information to make main.wasm smaller
RUN tinygo build -o /site/main.wasm -target wasm -no-debug .

# wasm_exec.js is the JavaScript helper that starts main.wasm in the browser.
# TinyGo has its OWN version of this file - the one from standard Go does not work here.
RUN cp "$(tinygo env TINYGOROOT)/targets/wasm_exec.js" /site/

# Add our web page (index.html, app.js, style.css).
COPY web/ /site/

# ---------------------------------------------------------------
# Stage 2: export (optional)
# ---------------------------------------------------------------
# An empty image that holds only the two built files.
# Docker can write them into the web/ folder on your disk:
#
#   docker build --target export --output web .
#
FROM scratch AS export
COPY --from=build /site/main.wasm /site/wasm_exec.js /

# ---------------------------------------------------------------
# Stage 3: serve
# ---------------------------------------------------------------
FROM nginx:1.30.6-alpine AS serve

# nginx serves everything in this folder.
# It already sends .wasm files with the "application/wasm" content type,
# which the browser requires for WebAssembly.instantiateStreaming().
COPY --from=build /site/ /usr/share/nginx/html/
