// This file is the whole Go side of the calculator.
// It is compiled to WebAssembly (main.wasm) by TinyGo and runs inside the browser.

// The line below tells Go to compile this file only for the browser
// (operating system "js", architecture "wasm").
//go:build js && wasm

package main

import (
	"errors"
	"syscall/js" // Go's package for talking to JavaScript
)

// Divide is the actual calculation.
// It is plain Go code with nothing WebAssembly-specific in it.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

// divide is the bridge between JavaScript and Go.
//
// JavaScript calls it as divide(a, b). Every JavaScript argument
// arrives as a js.Value, so we convert both to Go numbers with Float().
//
// The returned Go map becomes a plain JavaScript object:
//
//	{ result: 2.5 }                       when the division works
//	{ error: "cannot divide by zero" }    when b is zero
func divide(this js.Value, args []js.Value) any {
	a := args[0].Float()
	b := args[1].Float()

	result, err := Divide(a, b)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	return map[string]any{"result": result}
}

func main() {
	// Put our Go function on JavaScript's global object (window),
	// so the web page can call it as divide(a, b).
	js.Global().Set("divide", js.FuncOf(divide))

	// Wait forever. If main() returned, the Go program would stop
	// and JavaScript could no longer call divide.
	select {}
}
