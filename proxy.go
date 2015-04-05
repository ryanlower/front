package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

type Proxy struct {
	config           Config
	contentTypeRegex *regexp.Regexp
}

func newProxy(config Config) *Proxy {
	// Todo, add config error checking

	contentTypeRegex, err := regexp.Compile(config.allowedContentTypes)
	if err != nil {
		log.Panic(err)
	}

	return &Proxy{
		config:           config,
		contentTypeRegex: contentTypeRegex,
	}
}

func (p Proxy) handler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	if err := p.validRequest(params); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		p.proxyRequest(w, params)
	}
}

func (p Proxy) validRequest(params url.Values) error {
	url := params.Get("url")
	if url == "" {
		return errors.New("No request url to proxy")
	}

	return nil
}

func (p Proxy) proxyRequest(w http.ResponseWriter, params url.Values) {
	url := params.Get("url")

	resp, err := http.Get(url) // http.Get follows up to 10 redirects

	defer resp.Body.Close()

	if err != nil {
		log.Print(err)
		// Todo, handle specific errors
		http.Error(w, "Could not proxy", http.StatusInternalServerError)
	}

	if resp.StatusCode != 200 {
		log.Printf("Upstream response: %v", resp.StatusCode)
		http.Error(w, "Could not proxy", resp.StatusCode)
		return
	}

	if err := p.validResponse(resp); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	width := params.Get("width")
	height := params.Get("height")
	if width != "" && height != "" {
		img, _ := newImg(resp.Body)
		if err != nil {
			// TODO, handle error (redirect to original url?)
			log.Println(err)
			http.Error(w, "Could not create img", http.StatusInternalServerError)
			return
		}

		// TODO, handle missing query params
		widthInt, _ := strconv.Atoi(width)
		heightInt, _ := strconv.Atoi(height)
		img.resize(widthInt, heightInt)

		// Write the resized image (as jpeg) to the http.ResponseWriter
		img.write(w)
	} else {
		// Just copy the upstream response to the http.ResponseWriter
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		io.Copy(w, resp.Body)
	}
}

func (p Proxy) validResponse(resp *http.Response) error {
	if !p.contentTypeRegex.MatchString(resp.Header.Get("Content-Type")) {
		return errors.New("Upstream content doesn't match configured allowedContentTypes")
	}

	return nil
}
