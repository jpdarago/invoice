package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestGetImageDimensionMissing(t *testing.T) {
	w, h, err := getImageDimension("/nonexistent/path/image.png")
	if err == nil {
		t.Error("expected error for missing file")
	}
	if w != 0 || h != 0 {
		t.Errorf("got (%d, %d), want (0, 0) for missing file", w, h)
	}
}

func TestGetImageDimensionValid(t *testing.T) {
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "test.png")

	img := image.NewRGBA(image.Rect(0, 0, 80, 40))
	img.Set(0, 0, color.White)

	f, err := os.Create(imgPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()

	w, h, err2 := getImageDimension(imgPath)
	if err2 != nil {
		t.Fatalf("unexpected error: %v", err2)
	}
	if w != 80 || h != 40 {
		t.Errorf("got (%d, %d), want (80, 40)", w, h)
	}
}
