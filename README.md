# kirkmeme

A CLI tool that overlays meme-style text onto the Kirk image.

## Setup

```bash
go mod tidy
```

## Usage

```bash
go run main.go "KHAAAAN!" output.png
```

Or build and run:

```bash
go build -o kirkmeme .
./kirkmeme "KHAAAAN!"
./kirkmeme "He tasks me" rage.png
./kirkmeme --color yellow "KHAAAAN!"
./kirkmeme --color "#ff4500" "He tasks me" rage.png
```

- Text is automatically uppercased.
- Output defaults to `output.png` if not specified.
- `kirk.png` is embedded in the binary — no external files needed.

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `--color` | `white` | Text fill color — a name or hex value |

**Named colors:** `white`, `black`, `red`, `green`, `blue`, `yellow`, `orange`, `cyan`, `magenta`, `pink`

**Hex colors:** `#RRGGBB` or `#RGB` (e.g. `#ff4500`)

## Dependencies

- [fogleman/gg](https://github.com/fogleman/gg) — 2D graphics library
- A bold TTF font installed on your system (DejaVu, Liberation, Arial, Impact, etc.)
