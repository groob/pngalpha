package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"testing"
)

func TestConvert(t *testing.T) {
	file := newFile(t)
	img, err := convert(file, "image/jpeg")
	if err != nil {
		t.Fatal(err)
	}
	type opaquer interface {
		Opaque() bool
	}

	o := img.(opaquer)
	if o.Opaque() {
		t.Error("expected resulting image to not be opaque")
	}
}

// create a io.Reader to use in tests
func newFile(t *testing.T) *bytes.Buffer {
	m := image.NewRGBA(image.Rect(0, 0, 640, 480))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	opaque := m.Opaque()
	if !opaque {
		t.Fatal("the test image needs to be opaque")
	}
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, m, nil)
	if err != nil {
		t.Fatal(err)
	}
	return &buf
}
