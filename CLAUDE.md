# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go mod tidy          # install/sync dependencies
go build -o kirkmeme .  # build binary
go run main.go "TEXT" [output.png]  # run without building
```

No tests exist in this project.

## Architecture

Single-file Go CLI (`main.go`) that:

1. Embeds `kirk.png` at compile time via `//go:embed`
2. Parses `--color` flag and positional args (text, optional output path)
3. Uses `github.com/fogleman/gg` (2D graphics via freetype) to draw text onto the image
4. Auto-sizes font via binary search to fit text within 90% width / 20% height of the image
5. Renders 8-direction black outline, then colored fill text, anchored to the bottom of the image
6. Saves result as PNG

Key functions: `wrapText` (manual word-wrap matching gg's layout), `parseColor` (name or `#RRGGBB`/`#RGB` hex), `defaultFontPaths` (cross-platform bold TTF discovery).

**Font requirement:** A bold TTF must exist on the system — the binary does not bundle a font. See `defaultFontPaths()` for the search order across Linux/macOS/Windows.
