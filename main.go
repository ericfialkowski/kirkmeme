package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

//go:embed kirk.png
var kirkPNG []byte

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <text> [output.png]\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\nOverlays TEXT onto kirk.png in meme style.\n")
		fmt.Fprintf(os.Stderr, "Output defaults to output.png if not specified.\n")
		os.Exit(1)
	}

	text := strings.ToUpper(os.Args[1])
	outPath := "output.png"
	if len(os.Args) >= 3 {
		outPath = os.Args[2]
	}

	img, _, err := image.Decode(bytes.NewReader(kirkPNG))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading embedded image: %v\n", err)
		os.Exit(1)
	}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	dc := gg.NewContext(w, h)
	dc.DrawImage(img, 0, 0)

	// Find the first available bold system font that freetype can load.
	fontPath := ""
	for _, p := range defaultFontPaths() {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		if err := dc.LoadFontFace(p, 12); err == nil {
			fontPath = p
			break
		}
	}
	if fontPath == "" {
		fmt.Fprintln(os.Stderr, "error: no suitable bold TTF font found on this system")
		os.Exit(1)
	}

	// Horizontal padding: text must fit within 90% of image width.
	maxWidth := float64(w) * 0.9
	// Vertical budget: text must not exceed 20% of image height.
	maxHeight := float64(h) * 0.20

	// Binary-search for the largest font size whose wrapped text fits
	// within both maxWidth and maxHeight.
	lo, hi := 10.0, float64(w)
	fontSize := lo
	for lo <= hi {
		mid := math.Floor((lo + hi) / 2)
		if err := dc.LoadFontFace(fontPath, mid); err != nil {
			break
		}
		textW, textH := dc.MeasureMultilineString(wrapText(dc, text, maxWidth), 1.2)
		if textH <= maxHeight && textW <= maxWidth {
			fontSize = mid // fits — try bigger
			lo = mid + 1
		} else {
			hi = mid - 1 // too tall — try smaller
		}
	}

	// Reload the winning font size.
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		fmt.Fprintf(os.Stderr, "error loading font: %v\n", err)
		os.Exit(1)
	}

	// Pre-wrap once so measurement and drawing are identical.
	wrapped := wrapText(dc, text, maxWidth)
	_, textH := dc.MeasureMultilineString(wrapped, 1.2)

	// Anchor the TOP of the text block so the bottom sits near the image bottom.
	x := float64(w) / 2
	yTop := float64(h) - fontSize*0.2 - textH

	// Draw black outline by rendering text offset in 8 directions.
	dc.SetRGB(0, 0, 0)
	outline := fontSize / 15
	if outline < 1 {
		outline = 1
	}
	for dx := -outline; dx <= outline; dx += outline {
		for dy := -outline; dy <= outline; dy += outline {
			if dx == 0 && dy == 0 {
				continue
			}
			dc.DrawStringWrapped(wrapped, x+dx, yTop+dy, 0.5, 0, maxWidth, 1.2, gg.AlignCenter)
		}
	}

	// Draw white fill on top.
	dc.SetRGB(1, 1, 1)
	dc.DrawStringWrapped(wrapped, x, yTop, 0.5, 0, maxWidth, 1.2, gg.AlignCenter)

	if err := dc.SavePNG(outPath); err != nil {
		fmt.Fprintf(os.Stderr, "error saving image: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("saved to %s\n", outPath)
}

// wrapText simulates gg's word-wrapping to produce the multi-line string
// that MeasureMultilineString needs for accurate height measurement.
func wrapText(dc *gg.Context, text string, maxWidth float64) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}
	var lines []string
	cur := words[0]
	for _, w := range words[1:] {
		test := cur + " " + w
		tw, _ := dc.MeasureString(test)
		if tw > maxWidth {
			lines = append(lines, cur)
			cur = w
		} else {
			cur = test
		}
	}
	lines = append(lines, cur)
	return strings.Join(lines, "\n")
}

// defaultFontPaths returns common locations for a bold TTF on Linux / macOS / Windows.
func defaultFontPaths() []string {
	return []string{
		// Linux
		"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
		"/usr/share/fonts/truetype/ubuntu/Ubuntu-B.ttf",
		"/usr/share/fonts/TTF/DejaVuSans-Bold.ttf",
		// macOS
		"/System/Library/Fonts/Supplemental/Impact.ttf",
		"/Library/Fonts/Arial Bold.ttf",
		"/System/Library/Fonts/Helvetica.ttc",
		"/System/Library/Fonts/SFCompact.ttf",
		// Windows
		`C:\Windows\Fonts\arialbd.ttf`,
		`C:\Windows\Fonts\impact.ttf`,
	}
}
