package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Parse command-line arguments
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input-file>\n", os.Args[0])
		os.Exit(1)
	}

	inputPath := os.Args[1]

	// Read input file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %v\n", inputPath, err)
		os.Exit(1)
	}

	// Detect encoding
	encoding, err := DetectEncoding(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting encoding: %v\n", err)
		os.Exit(3)
	}

	fmt.Printf("Detected encoding: %s\n", encoding)

	// Check if already UTF-8
	if encoding == "utf-8" {
		fmt.Println("File is already in UTF-8 format. No conversion needed.")
		os.Exit(0)
	}

	// Convert to UTF-8
	utf8Data, err := ConvertToUTF8(data, encoding)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting to UTF-8: %v\n", err)
		os.Exit(3)
	}

	// Generate output filename
	outputPath := generateOutputFilename(inputPath)

	// Write output file
	err = os.WriteFile(outputPath, utf8Data, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file '%s': %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted '%s' to UTF-8\n", inputPath)
	fmt.Printf("Output written to: %s\n", outputPath)
}

// generateOutputFilename inserts "-utf8" before the file extension
// Examples:
//   - "test.txt" -> "test-utf8.txt"
//   - "file" -> "file-utf8"
//   - "archive.tar.gz" -> "archive.tar-utf8.gz"
func generateOutputFilename(inputPath string) string {
	dir := filepath.Dir(inputPath)
	filename := filepath.Base(inputPath)

	// Find the last dot to split name and extension
	ext := filepath.Ext(filename)

	if ext == "" {
		// No extension
		return filepath.Join(dir, filename+"-utf8")
	}

	// Split filename into base name and extension
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// Generate new filename
	newFilename := nameWithoutExt + "-utf8" + ext

	return filepath.Join(dir, newFilename)
}
