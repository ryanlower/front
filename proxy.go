package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Proxy struct {
	config           *config
	contentTypeRegex *regexp.Regexp
	s3               *s3
}

func newProxy(config *config, s3 *s3) *Proxy {
	// Todo, add config error checking

	contentTypeRegex, err := regexp.Compile(config.AllowedContentTypes)
	if err != nil {
		log.Panic(err)
	}

	return &Proxy{
		config:           config,
		contentTypeRegex: contentTypeRegex,
		s3:               s3,
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
	path := params.Get("path")
	if url == "" && path == "" {
		return errors.New("No request url or path to proxy")
	}

	return nil
}

func (p Proxy) proxyRequest(w http.ResponseWriter, params url.Values) {
	// TODO, handle missing (or zero) query params
	widthInt, _ := strconv.Atoi(params.Get("width"))
	heightInt, _ := strconv.Atoi(params.Get("height"))

	if path := params.Get("path"); path != "" {
		cachedPath := strings.ToLower(
			fmt.Sprintf("front/%dx%d/%s", widthInt, heightInt, path))

		if body, etag, err := p.s3.read(cachedPath); err == nil {
			// S3 cache hit, return cached body
			log.Printf("[HIT] %s", cachedPath)
			w.Header().Set("Etag", etag)

			io.Copy(w, body)
		} else {
			// S3 cache miss, get path from S3, resize, return (and write to cache)
			log.Printf("[MISS] %s", cachedPath)
			body, etag, err := p.s3.read(path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			defer body.Close()

			img, err := newImg(body)
			if err != nil {
				s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", p.config.S3.Bucket, path)

				log.Printf("[REDIRECT] %s", s3URL)
				w.Header().Set("Location", s3URL)
				w.WriteHeader(http.StatusFound)
				return
			}
			img.resize(widthInt, heightInt)

			w.Header().Set("Etag", etag)
			// Write the resized image (as jpeg) to the http.ResponseWriter
			img.write(w)
			// Write the resized image to s3 cache
			// TODO, handle writing errors?
			go p.s3.write(cachedPath, img)
		}
	} else if url := params.Get("url"); url != "" {
		resp, err := http.Get(url) // http.Get follows up to 10 redirects
		if err != nil {
			http.Error(w, "Could not proxy", http.StatusInternalServerError)
			return
		}
		if resp.StatusCode != 200 {
			http.Error(w, "Could not proxy", resp.StatusCode)
			return
		}
		defer resp.Body.Close()

		if err := p.validResponse(resp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if widthInt > 0 && heightInt > 0 {
			img, err := newImg(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			img.resize(widthInt, heightInt)
			// Write the resized image (as jpeg) to the http.ResponseWriter
			img.write(w)
		} else {
			// Just copy the upstream response to the http.ResponseWriter
			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			io.Copy(w, resp.Body)
		}
	}
}

func (p Proxy) validResponse(resp *http.Response) error {
	if !p.contentTypeRegex.MatchString(resp.Header.Get("Content-Type")) {
		return errors.New("Upstream content doesn't match configured allowedContentTypes")
	}

	return nil
}
