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
	w, h := getImageDimension("/nonexistent/path/image.png")
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

	w, h := getImageDimension(imgPath)
	if w != 80 || h != 40 {
		t.Errorf("got (%d, %d), want (80, 40)", w, h)
	}
}
