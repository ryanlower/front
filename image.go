package main

import (
	"io/ioutil"
	"net/http"
)

type Image struct {
	bytes []byte
}

func NewImage(bytes []byte) (image *Image) {
	return &Image{
		bytes: bytes,
	}
}

func NewImageFromResponse(resp *http.Response) (image *Image) {
	// TODO, error handling
	body, _ := ioutil.ReadAll(resp.Body)
	return NewImage(body)
}

func (i *Image) Read() []byte {
	return i.bytes
}
