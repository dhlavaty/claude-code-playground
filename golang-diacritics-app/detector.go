package main

import (
	"bytes"
	"errors"
)

// DetectEncoding detects whether the input data is UTF-8, WINDOWS-1250, or ISO-8859-2
func DetectEncoding(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empty input data")
	}

	// Check for UTF-8 BOM (0xEF 0xBB 0xBF)
	if len(data) >= 3 && bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF}) {
		return "utf-8", nil
	}

	// Count markers for WINDOWS-1250 and ISO-8859-2
	windows1250Markers := 0
	iso88592Markers := 0

	// WINDOWS-1250 specific bytes in the 0x80-0x9F range (undefined in ISO-8859-2)
	// These are common in Czech/Slovak text: Š=0x8A, Ť=0x8D, Ž=0x8E, š=0x9A, ť=0x9D, ž=0x9E
	windows1250Specific := []byte{0x8A, 0x8D, 0x8E, 0x9A, 0x9D, 0x9E}

	// ISO-8859-2 specific bytes that differ from WINDOWS-1250
	// Key distinguishing characters in ISO-8859-2:
	// - Ľ=0xA5 (in WIN1250 it's Ą=0xA5)
	// - Ś=0xA6 (different positions)
	// - Ť=0xAB (in WIN1250 it's at 0x8D)
	// - ľ=0xB5 (in WIN1250 it's µ=0xB5)
	// - ś=0xB6 (different positions)
	// - ť=0xBB (in WIN1250 it's at 0x9D)
	// - ź=0xBC (different positions)
	iso88592Specific := []byte{0xA5, 0xA6, 0xAB, 0xB5, 0xB6, 0xBB, 0xBC}

	// Scan the data for encoding markers
	for _, b := range data {
		// Check for WINDOWS-1250 markers
		for _, marker := range windows1250Specific {
			if b == marker {
				windows1250Markers++
				break
			}
		}

		// Check for ISO-8859-2 markers
		for _, marker := range iso88592Specific {
			if b == marker {
				iso88592Markers++
				break
			}
		}
	}

	// Make a decision based on marker counts
	if windows1250Markers > 0 && windows1250Markers > iso88592Markers {
		return "windows-1250", nil
	}

	if iso88592Markers > 0 && iso88592Markers > windows1250Markers {
		return "iso-8859-2", nil
	}

	// If we have both markers equally or no clear winner, check more carefully
	if windows1250Markers > 0 && iso88592Markers > 0 {
		// If there are markers in the 0x80-0x9F range, it's definitely WINDOWS-1250
		// because ISO-8859-2 doesn't define these
		for _, b := range data {
			if b >= 0x80 && b <= 0x9F {
				for _, marker := range windows1250Specific {
					if b == marker {
						return "windows-1250", nil
					}
				}
			}
		}
	}

	// If no distinguishing markers found, check if it's pure ASCII
	isASCII := true
	for _, b := range data {
		if b >= 0x80 {
			isASCII = false
			break
		}
	}

	if isASCII {
		// Pure ASCII is compatible with all encodings, default to WINDOWS-1250
		return "windows-1250", nil
	}

	// Default to WINDOWS-1250 (more common for Czech/Slovak texts)
	return "windows-1250", nil
}
