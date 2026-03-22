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
```

- Text is automatically uppercased.
- Output defaults to `output.png` if not specified.
- Keep `kirk.png` in the same directory as the binary (or cwd).

## Dependencies

- [fogleman/gg](https://github.com/fogleman/gg) — 2D graphics library
- A bold TTF font installed on your system (DejaVu, Liberation, Arial, Impact, etc.)
