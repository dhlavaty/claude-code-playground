# Diacritics converter application

## Functional Specification Document (FSD)

**Project:** Application for fully automatic conversion of text files from WINDOWS-1250 or ISO-8859-2 to UTF-8

## Architecture

- Command-line application written in GOlang (The Go Programming Language)
- Application will use as little dependencies (external libraries) as possible
- Project contains simple helper shell scripts
  - `run.sh` will compile and run the application
  - `build.sh` will compile the application

## Purpose

Application serves for automatic conversion of text files from WINDOWS-1250 or ISO-8859-2 to UTF-8 encoding with BOM (Byte order mark). Application will use heuristic to detect if input file is in WINDOWS-1250 or ISO-8859-2 encoding.

## Input Handling

Application is a command-line application, where the first argument is the input file. Output file is saved with the original name, prepended with `-utf8` text.

## Testing files

There are testing files in `test-files/` sub folder.

File [test1-windows1250.txt](golang-diacritics-app/test-files/test1-windows1250.txt) contais an example text in Slovak language encoded in WINDOWS-1250 encoding. File [test1-iso8859-2.txt](golang-diacritics-app/test-files/test1-iso8859-2.txt) contains the same text encoded using ISO-8859-2 encoding. There is also a [test1-utf8.txt](golang-diacritics-app/test-files/test1-utf8.txt) file, with the same text encoded in UTF-8 with BOM (Byte order mark).
