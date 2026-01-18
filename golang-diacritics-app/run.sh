#!/bin/bash
set -e
if [ $# -eq 0 ]; then
    echo "Usage: ./run.sh <input-file>"
    exit 1
fi
echo "Building..."
~/go/bin/go build -o diacritics-converter main.go detector.go converter.go
echo "Running conversion on: $1"
./diacritics-converter "$1"
echo "Conversion complete!"
