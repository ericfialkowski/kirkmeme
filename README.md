# kirkmeme

A CLI tool that overlays meme-style text onto the Kirk image.

## Installation

### Homebrew (macOS)

```bash
brew install --cask ericfialkowski/tap/kirkmeme
```

### Download a release

Download the archive for your platform from the [Releases page](https://github.com/ericfialkowski/kirkmeme/releases), extract it, and place the `kirkmeme` binary somewhere on your `PATH`.

### Build from source

Requires Go 1.24+ and a bold TTF font on your system (see [Dependencies](#dependencies)).

```bash
go install github.com/ericfialkowski/kirkmeme@latest
```

Or clone and build:

```bash
git clone https://github.com/ericfialkowski/kirkmeme.git
cd kirkmeme
go build -o kirkmeme .
```

## Usage

```bash
kirkmeme "KHAAAAN!"
kirkmeme "He tasks me" rage.png
kirkmeme --color yellow "KHAAAAN!"
kirkmeme --color "#ff4500" "He tasks me" rage.png
kirkmeme --no-exclaim --file "He's dead, Jim"
```

- Text is automatically uppercased.
- Exclamation marks are added automatically if missing (disable with `--no-exclaim`).
- The result is copied to the clipboard by default (disable with `--clipboard=false`).
- Output file defaults to `output.png` if not specified.
- `kirk.png` is embedded in the binary — no external files needed.

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `--color` | `white` | Text fill color — a name or hex value |
| `--clipboard` | `true` | Copy result to system clipboard |
| `--file` | `false` | Write result to a file |
| `--no-exclaim` | `false` | Disable auto-adding exclamation marks |

**Named colors:** `white`, `black`, `red`, `green`, `blue`, `yellow`, `orange`, `cyan`, `magenta`, `pink`

**Hex colors:** `#RRGGBB` or `#RGB` (e.g. `#ff4500`)

## Dependencies

- [fogleman/gg](https://github.com/fogleman/gg) — 2D graphics library
- [golang.design/x/clipboard](https://pkg.go.dev/golang.design/x/clipboard) — clipboard access (requires CGO; Linux also needs X11 dev headers)
- A bold TTF font installed on your system (DejaVu, Liberation, Arial, Impact, etc.)
