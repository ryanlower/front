package main

import (
	"image"
	"image/jpeg"
	"io"
	"net/http"

	"github.com/disintegration/imaging"
)

type img struct {
	image image.Image
}

// new creates a new img from r
func newImg(r io.Reader) (*img, error) {
	i, err := imaging.Decode(r)
	if err != nil {
		return nil, err
	}

	return &img{image: i}, nil
}

// adds content-type header and encodes image to w
func (i *img) write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", i.contentType())
	i.encode(w)
}

func (i *img) encode(w io.Writer) {
	// TODO, how to handle formats?
	// Should we always encode as X (jpeg?),
	// or encode in the same format as the original image

	// default quality encoding is 75
	jpeg.Encode(w, i.image, nil)
}

func (i *img) resize(width, height int) {
	// TODO, allow customisation of filter
	filter := imaging.Linear

	i.image = imaging.Thumbnail(i.image, width, height, filter)
}

func (i *img) contentType() string {
	// Always jpeg at the moment (see encode)
	return "image/jpeg"
}
