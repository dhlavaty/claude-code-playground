// Find the page elements we work with.
const form = document.getElementById("calculator");
const inputA = document.getElementById("a");
const inputB = document.getElementById("b");
const button = document.getElementById("divide");
const output = document.getElementById("output");
const status = document.getElementById("status");

// ---------------------------------------------------------------
// Step 1: load and start the WebAssembly module.
// ---------------------------------------------------------------

// Go comes from wasm_exec.js. It provides the functions
// the Go program needs from the browser (go.importObject).
const go = new Go();

WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
  .then((wasm) => {
    // Start the Go program. Its main() puts divide() on window.
    go.run(wasm.instance);

    button.disabled = false;
    status.textContent = "WebAssembly loaded. Ready.";
  })
  .catch((error) => {
    status.textContent = "Could not load WebAssembly: " + error;
  });

// ---------------------------------------------------------------
// Step 2: when the form is submitted, call the Go function.
// ---------------------------------------------------------------

form.addEventListener("submit", (event) => {
  event.preventDefault(); // stay on the page, don't reload it

  const a = toNumber(inputA.value);
  const b = toNumber(inputB.value);

  // Checking the input is the page's job.
  // The division rule (no dividing by zero) is Go's job.
  if (Number.isNaN(a) || Number.isNaN(b)) {
    showError("Please enter two numbers.");
    return;
  }

  // This calls Go! The answer is { result: ... } or { error: ... }.
  const answer = divide(a, b);

  if (answer.error) {
    showError(answer.error);
  } else {
    showResult(a + " ÷ " + b + " = " + answer.result);
  }
});

// Turns the text from an input into a number.
// Accepts "," as a decimal separator too, e.g. "2,5".
// Returns NaN ("not a number") for empty or invalid text.
function toNumber(text) {
  const cleaned = text.trim().replace(",", ".");
  if (cleaned === "") {
    return NaN;
  }
  return Number(cleaned);
}

function showResult(text) {
  output.textContent = text;
  output.classList.remove("error");
}

function showError(text) {
  output.textContent = text;
  output.classList.add("error");
}
