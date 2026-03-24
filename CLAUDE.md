# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go mod tidy              # install/sync dependencies
go build -o kirkmeme .   # build binary
go run main.go "TEXT" [output.png]  # run without building
go test ./...            # run tests
```

## Architecture

Single-file Go CLI (`main.go`) that:

1. Embeds `kirk.png` at compile time via `//go:embed`
2. Parses flags (`--color`, `--clipboard`, `--file`, `--no-exclaim`) and positional args (text, optional output path)
3. Auto-appends 1–3 exclamation marks if the text doesn't already end with one (suppressible)
4. Uses `github.com/fogleman/gg` (2D graphics via freetype) to draw text onto the image
5. Auto-sizes font via binary search to fit text within 90% width / 20% height of the image
6. Renders 8-direction black outline, then colored fill text, anchored to the bottom of the image
7. Encodes result as PNG; copies to clipboard and/or saves to file based on flags

Key functions: `wrapText` (manual word-wrap matching gg's layout), `parseColor` (name or `#RRGGBB`/`#RGB` hex), `defaultFontPaths` (cross-platform bold TTF discovery).

**Font requirement:** A bold TTF must exist on the system — the binary does not bundle a font. See `defaultFontPaths()` for the search order across Linux/macOS/Windows.

**CGO requirement:** `golang.design/x/clipboard` requires CGO. On Linux, X11 dev headers (`libx11-dev`, `libxcb1-dev`) must be present. Windows builds use `CGO_ENABLED=0` (clipboard uses pure Go syscalls there).

## Release workflow

Releases are triggered by pushing a `v*` tag. The workflow uses three native OS runners so CGO compiles against the correct platform toolchain:

| Runner | Config | Targets |
|--------|--------|---------|
| `ubuntu-latest` | `.goreleaser.linux.yml` | linux/amd64 |
| `macos-latest` | `.goreleaser.darwin.yml` | darwin/amd64, darwin/arm64 |
| `windows-latest` | `.goreleaser.windows.yml` | windows/amd64, windows/arm64 |

Each runner builds and archives with `goreleaser release --skip=publish,announce,validate`. A final `release` job collects all artifacts, publishes the GitHub release via `gh`, and updates the Homebrew tap.
