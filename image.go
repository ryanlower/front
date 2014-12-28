package main

import (
	"image"
	"net/http"

	"github.com/disintegration/imaging"
)

type Image struct {
	img    image.Image
	format imaging.Format
}

func NewImageFromResponse(resp *http.Response) *Image {
	// TODO, error handling
	img, formatString, _ := image.Decode(resp.Body)
	format := stringToFormat(formatString)

	return &Image{
		img:    img,
		format: format,
	}
}

func stringToFormat(s string) imaging.Format {
	// These are the formats supported by the imaging package
	formats := map[string]imaging.Format{
		"jpeg": imaging.JPEG,
		"png":  imaging.PNG,
		"gif":  imaging.GIF,
		"tiff": imaging.TIFF,
		"bmp":  imaging.BMP,
	}

	return formats[s]
}

func (i *Image) Scale(factor float32) {
	width, height := i.getSize()
	newWidth := float32(width) * factor
	newHeight := float32(height) * factor
	i.resize(int(newWidth), int(newHeight))
}

func (i *Image) Write(w http.ResponseWriter) {
	imaging.Encode(w, i.img, i.format)
}

func (i *Image) getSize() (int, int) {
	bounds := i.img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

func (i *Image) resize(width, height int) {
	i.img = imaging.Resize(i.img, width, height, imaging.Box)
}
