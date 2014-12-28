package main

import (
	"errors"
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
	if err != nil {
		log.Print(err)
		// Todo, handle specific errors
		http.Error(w, "Could not proxy", http.StatusInternalServerError)
	}
	if resp.StatusCode != 200 {
		log.Printf("Upstream response: %v", resp.StatusCode)
		http.Error(w, "Could not proxy", resp.StatusCode)
	}

	if err := p.validResponse(resp); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	p.writeResponse(w, resp, params)
}

func (p Proxy) validResponse(resp *http.Response) error {
	if !p.contentTypeRegex.MatchString(resp.Header.Get("Content-Type")) {
		return errors.New("Upstream content doesn't match configured allowedContentTypes")
	}

	return nil
}

func (p Proxy) writeResponse(w http.ResponseWriter, resp *http.Response, params url.Values) {
	defer resp.Body.Close()

	// TODO, handle non images (just return upstream content as is?)
	image := NewImageFromResponse(resp)

	// Scale image if scale param is present
	if scale := params.Get("scale"); scale != "" {
		scaleFloat, _ := strconv.ParseFloat(scale, 32)
		image.Scale(float32(scaleFloat))
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	image.Write(w)
}
