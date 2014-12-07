package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
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

	if p.contentTypeRegex.MatchString(resp.Header.Get("Content-Type")) {
		p.writeResponse(w, resp)
	} else {
		log.Println("Upstream content doesn't match configured allowedContentTypes")
		http.Error(w, "Upstream content doesn't match configured allowedContentTypes", http.StatusBadRequest)
	}
}

func (p Proxy) writeResponse(w http.ResponseWriter, resp *http.Response) {
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	io.Copy(w, resp.Body)
	resp.Body.Close()
}
