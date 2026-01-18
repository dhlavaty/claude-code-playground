#!/bin/bash
set -e
echo "Building diacritics converter..."
~/go/bin/go build -o diacritics-converter main.go detector.go converter.go
echo "Build successful! Binary: diacritics-converter"
