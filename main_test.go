package main

import (
	"math"
	"testing"

	"github.com/fogleman/gg"
)

// ---- parseColor tests ----

func TestParseColorNamed(t *testing.T) {
	tests := []struct {
		name string
		r, g, b float64
	}{
		{"white", 1, 1, 1},
		{"black", 0, 0, 0},
		{"red", 1, 0, 0},
		{"yellow", 1, 1, 0},
		{"cyan", 0, 1, 1},
		{"magenta", 1, 0, 1},
	}
	for _, tt := range tests {
		r, g, b, err := parseColor(tt.name)
		if err != nil {
			t.Errorf("parseColor(%q) unexpected error: %v", tt.name, err)
			continue
		}
		if r != tt.r || g != tt.g || b != tt.b {
			t.Errorf("parseColor(%q) = (%.2f, %.2f, %.2f), want (%.2f, %.2f, %.2f)",
				tt.name, r, g, b, tt.r, tt.g, tt.b)
		}
	}
}

func TestParseColorNamedCaseInsensitive(t *testing.T) {
	for _, s := range []string{"WHITE", "White", "wHiTe"} {
		r, g, b, err := parseColor(s)
		if err != nil {
			t.Errorf("parseColor(%q) unexpected error: %v", s, err)
			continue
		}
		if r != 1 || g != 1 || b != 1 {
			t.Errorf("parseColor(%q) = (%.2f, %.2f, %.2f), want (1, 1, 1)", s, r, g, b)
		}
	}
}

func TestParseColorHexRRGGBB(t *testing.T) {
	tests := []struct {
		input   string
		r, g, b float64
	}{
		{"#ff0000", 1, 0, 0},
		{"#00ff00", 0, 1, 0},
		{"#0000ff", 0, 0, 1},
		{"#ffffff", 1, 1, 1},
		{"#000000", 0, 0, 0},
		{"#804020", float64(0x80) / 255, float64(0x40) / 255, float64(0x20) / 255},
	}
	for _, tt := range tests {
		r, g, b, err := parseColor(tt.input)
		if err != nil {
			t.Errorf("parseColor(%q) unexpected error: %v", tt.input, err)
			continue
		}
		eps := 1e-9
		if math.Abs(r-tt.r) > eps || math.Abs(g-tt.g) > eps || math.Abs(b-tt.b) > eps {
			t.Errorf("parseColor(%q) = (%.6f, %.6f, %.6f), want (%.6f, %.6f, %.6f)",
				tt.input, r, g, b, tt.r, tt.g, tt.b)
		}
	}
}

func TestParseColorHexRGB(t *testing.T) {
	// #RGB should expand to #RRGGBB
	r, g, b, err := parseColor("#f00")
	if err != nil {
		t.Fatalf("parseColor(\"#f00\") unexpected error: %v", err)
	}
	if r != 1 || g != 0 || b != 0 {
		t.Errorf("parseColor(\"#f00\") = (%.2f, %.2f, %.2f), want (1, 0, 0)", r, g, b)
	}

	r, g, b, err = parseColor("#fff")
	if err != nil {
		t.Fatalf("parseColor(\"#fff\") unexpected error: %v", err)
	}
	if r != 1 || g != 1 || b != 1 {
		t.Errorf("parseColor(\"#fff\") = (%.2f, %.2f, %.2f), want (1, 1, 1)", r, g, b)
	}
}

func TestParseColorErrors(t *testing.T) {
	bad := []string{
		"notacolor",
		"#gg0000",
		"#12345",   // 5 hex digits
		"#1234567", // 7 hex digits
		"",
	}
	for _, s := range bad {
		_, _, _, err := parseColor(s)
		if err == nil {
			t.Errorf("parseColor(%q) expected error, got nil", s)
		}
	}
}

// ---- defaultFontPaths tests ----

func TestDefaultFontPathsNonEmpty(t *testing.T) {
	paths := defaultFontPaths()
	if len(paths) == 0 {
		t.Error("defaultFontPaths() returned empty slice")
	}
}

func TestDefaultFontPathsNoDuplicates(t *testing.T) {
	seen := map[string]bool{}
	for _, p := range defaultFontPaths() {
		if seen[p] {
			t.Errorf("duplicate font path: %q", p)
		}
		seen[p] = true
	}
}

// ---- wrapText tests ----

func TestWrapTextEmpty(t *testing.T) {
	dc := gg.NewContext(100, 100)
	got := wrapText(dc, "", 500)
	if got != "" {
		t.Errorf("wrapText with empty text = %q, want %q", got, "")
	}
}

func TestWrapTextSingleWord(t *testing.T) {
	dc := gg.NewContext(100, 100)
	got := wrapText(dc, "HELLO", 500)
	if got != "HELLO" {
		t.Errorf("wrapText single word = %q, want %q", got, "HELLO")
	}
}

func TestWrapTextNoWrapNeeded(t *testing.T) {
	// With a very wide maxWidth no line breaks should appear.
	dc := gg.NewContext(100, 100)
	got := wrapText(dc, "ONE TWO THREE", 1e9)
	if got != "ONE TWO THREE" {
		t.Errorf("wrapText wide = %q, want %q", got, "ONE TWO THREE")
	}
}

func TestWrapTextForcesBreak(t *testing.T) {
	// With maxWidth=0 every word should land on its own line.
	dc := gg.NewContext(100, 100)
	got := wrapText(dc, "A B C", 0)
	want := "A\nB\nC"
	if got != want {
		t.Errorf("wrapText narrow = %q, want %q", got, want)
	}
}
