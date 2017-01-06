package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	var (
		flIn  = flag.String("in", "", "path to input image (jpeg or png)")
		flOut = flag.String("out", "/Library/Caches/com.apple.desktop.admin.png", "path to output png")
	)
	flag.Parse()

	in, err := os.Open(*flIn)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	ct, err := detectContentType(in)
	if err != nil {
		log.Fatal(err)
	}

	img, err := convert(in, ct)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(*flOut)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		log.Fatal(err)
	}
}

// detectContentType leverages the http.DetectContentType helper. It resets
// the io.ReadSeeker to 0 when done.
func detectContentType(file io.ReadSeeker) (string, error) {
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		return "", err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buf), nil
}

func convert(file io.Reader, contentType string) (image.Image, error) {
	var decode func(io.Reader) (image.Image, error)
	switch contentType {
	case "image/jpeg":
		decode = jpeg.Decode
	case "image/png":
		decode = png.Decode
	default:
		return nil, fmt.Errorf("unrecognized image format: %s", contentType)
	}
	src, err := decode(file)
	if err != nil {
		return nil, err
	}

	img := &notOpaqueRGBA{image.NewRGBA(src.Bounds())}
	draw.Draw(img, img.Bounds(), src, image.ZP, draw.Src)

	return img, nil
}

// enforce image.RGBA to always add the alpha channel when encoding PNGs.
type notOpaqueRGBA struct {
	*image.RGBA
}

func (i *notOpaqueRGBA) Opaque() bool {
	return false
}
