# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a single-file web application: a proportion calculator that solves for one missing value in the proportion equation `a/b = c/d`.

## Architecture

**Single-File Structure:**
- `index.html` - Complete self-contained application with embedded CSS and JavaScript
  - HTML: Fraction-style input layout with 4 fields (a, b, c, d)
  - CSS: Modern purple gradient design with responsive layout
  - JavaScript: Cross-multiplication solver with validation

**Core Calculation Logic (index.html:296-313):**
The calculator uses cross-multiplication to solve for the missing value:
- If `a` is missing: `a = (b × c) / d`
- If `b` is missing: `b = (a × d) / c`
- If `c` is missing: `c = (a × d) / b`
- If `d` is missing: `d = (b × c) / a`

**Validation Requirements:**
- Exactly 3 values must be provided (one field must be empty)
- All values must be positive numbers
- Division by zero is prevented with error handling

## Running the Application

Simply open `index.html` in a web browser:
```bash
xdg-open index.html  # Linux
open index.html      # macOS
start index.html     # Windows
```

Or use a local server:
```bash
python3 -m http.server 8000
# Then navigate to http://localhost:8000
```

## Making Changes

When modifying the application:
- All code is in one file - no build process required
- Test by refreshing the browser after making changes
- The application is fully client-side with no backend dependencies
