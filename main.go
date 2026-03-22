package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"golang.design/x/clipboard"
)

//go:embed kirk.png
var kirkPNG []byte

func main() {
	colorFlag := flag.String("color", "white", "text fill color: name (white, yellow, red, …) or hex (#RRGGBB)")
	clipboardFlag := flag.Bool("clipboard", true, "copy result to system clipboard")
	fileFlag := flag.Bool("file", false, "write result to a file")
	noExclaimFlag := flag.Bool("no-exclaim", false, "disable auto-adding exclamation marks to text that lacks them")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--color COLOR] [--clipboard=false] [--file] <text> [output.png]\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\nOverlays TEXT onto kirk.png in meme style.\n")
		fmt.Fprintf(os.Stderr, "Output file defaults to output.png if not specified.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	text := strings.ToUpper(flag.Arg(0))

	if !*noExclaimFlag && !strings.HasSuffix(text, "!") {
		n := rand.Intn(3) + 1
		text += strings.Repeat("!", n)
		fmt.Fprintf(os.Stderr, "Seriously? No exclamation mark? It's a MEME. Added %d for you.\n", n)
	}

	outPath := "output.png"
	if flag.NArg() >= 2 {
		outPath = flag.Arg(1)
	}

	fillR, fillG, fillB, err := parseColor(*colorFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
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

	// Draw fill on top in the chosen color.
	dc.SetRGB(fillR, fillG, fillB)
	dc.DrawStringWrapped(wrapped, x, yTop, 0.5, 0, maxWidth, 1.2, gg.AlignCenter)

	// Encode the result to a PNG byte buffer.
	var buf bytes.Buffer
	if err := dc.EncodePNG(&buf); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding image: %v\n", err)
		os.Exit(1)
	}
	pngBytes := buf.Bytes()

	if *clipboardFlag {
		if err := clipboard.Init(); err != nil {
			fmt.Fprintf(os.Stderr, "error initializing clipboard: %v\n", err)
			os.Exit(1)
		}
		clipboard.Write(clipboard.FmtImage, pngBytes)
		fmt.Println("copied to clipboard")
	}

	if *fileFlag {
		if err := os.WriteFile(outPath, pngBytes, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error saving image: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("saved to %s\n", outPath)
	}
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

// parseColor converts a color name or hex string to normalized [0,1] RGB components.
func parseColor(s string) (r, g, b float64, err error) {
	named := map[string][3]float64{
		"white":   {1, 1, 1},
		"black":   {0, 0, 0},
		"red":     {1, 0, 0},
		"green":   {0, 0.8, 0},
		"blue":    {0, 0.4, 1},
		"yellow":  {1, 1, 0},
		"orange":  {1, 0.5, 0},
		"cyan":    {0, 1, 1},
		"magenta": {1, 0, 1},
		"pink":    {1, 0.4, 0.7},
	}
	if c, ok := named[strings.ToLower(s)]; ok {
		return c[0], c[1], c[2], nil
	}
	hex := strings.TrimPrefix(s, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	if len(hex) != 6 {
		return 0, 0, 0, fmt.Errorf("unrecognized color %q (use a name or #RRGGBB)", s)
	}
	v, e := strconv.ParseUint(hex, 16, 32)
	if e != nil {
		return 0, 0, 0, fmt.Errorf("invalid hex color %q", s)
	}
	return float64(v>>16&0xff) / 255, float64(v>>8&0xff) / 255, float64(v&0xff) / 255, nil
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
