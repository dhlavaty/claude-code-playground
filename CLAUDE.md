# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a single-file web application: a proportion calculator that solves for one missing value in the proportion equation `a/b = c/d`.

## Architecture

**Single-File Structure:**
- `index.html` - Complete self-contained application with embedded CSS and JavaScript
  - HTML: Arrow-based input layout with 4 fields (a, b, c, d) displayed as `a → c` and `b → d`
  - CSS: Modern purple gradient design with responsive layout
  - JavaScript: Automatic cross-multiplication solver with real-time validation

**UI Layout:**
```
a → c
b → d
[Calculate] [Reset]
```

**Tab Order:**
The tabindex is configured for convenient data entry: a → c → b → d → Calculate → Reset

**Auto-Calculation Behavior:**
- Calculations happen automatically as you type - no need to press Calculate button
- When exactly 3 fields are filled, the 4th field is automatically calculated
- The calculated field is marked with a green border (`.calculated` class at index.html:93-97)
- As you continue typing or modify any value, the calculated field updates instantly
- If you type in the calculated field, it becomes a regular input and another field will be selected for calculation
- Errors (invalid input, division by zero) automatically clear the calculated field

**Input Handling (index.html:253-258):**
- Accepts both dot (.) and comma (,) as decimal separators
- Commas are automatically converted to dots before parsing: `.replace(/,/g, '.')`
- Supports European and American number formats
- Input event listeners trigger auto-calculation on every change (index.html:237-250)

**Core Calculation Logic (index.html:268-381):**
The calculator uses cross-multiplication to solve for the missing value:
- If `a` is missing: `a = (b × c) / d`
- If `b` is missing: `b = (a × d) / c`
- If `c` is missing: `c = (a × d) / b`
- If `d` is missing: `d = (b × c) / a`

Results are rounded to 2 decimal places: `Math.round(result * 100) / 100`

**Target Field Selection:**
- If exactly 1 field is empty, that field becomes the calculated target
- If a field is already calculated (has green border), it remains the target even when all 4 fields are filled
- User typing in a calculated field removes its calculated status

**Validation:**
- Exactly 3 input fields must have values (4th is auto-calculated)
- All values must be positive numbers
- Division by zero is prevented - clears calculated field
- Non-numeric input clears the calculated field
- All validation happens silently without error messages during typing

**Key Functions:**
- `getValues()` - Retrieves and parses input values, handles comma/dot conversion
- `autoCalculate()` - Main auto-calculation function triggered on every input change
- `calculate()` - Wrapper that calls autoCalculate() (for Calculate button compatibility)
- `reset()` - Clears all inputs, results, and calculated status
- `showResult(message, type)` - Displays success messages (errors are handled silently)

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
