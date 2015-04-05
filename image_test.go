package main

import (
	"image"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	assert := assert.New(t)
	w := httptest.NewRecorder()

	file, _ := os.Open("test/images/gopher.png")
	img, _ := newImg(file)

	img.write(w)

	assert.Equal(w.Code, 200, "status should be ok")
	assert.Equal(w.HeaderMap.Get("Content-Type"), "image/jpeg")

	file, _ = os.Open("test/images/gopher.png")
	fileImage, _, _ := image.Decode(file)
	bodyImage, _, _ := image.Decode(w.Body)
	assert.Equal(bodyImage.Bounds(), fileImage.Bounds(), "images should be same size")
}

func TestResize(t *testing.T) {
	assert := assert.New(t)

	file, _ := os.Open("test/images/gopher.png")
	img, _ := newImg(file)

	img.resize(100, 50)

	bounds := img.image.Bounds()
	assert.Equal(bounds.Dx(), 100, "width should be 100")
	assert.Equal(bounds.Dy(), 50, "height should be 50")
}
