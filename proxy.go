package main

import (
	"io"
	"log"
	"net/http"
)

func proxy(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("url")
	if request != "" {
		proxyRequest(w, request)
	} else {
		log.Println("No request url to proxy")
		http.Error(w, "No request url to proxy", http.StatusBadRequest)
	}
}

func proxyRequest(w http.ResponseWriter, request string) {
	resp, err := http.Get(request) // http.Get follows up to 10 redirects
	if err != nil {
		log.Print(err)
		// Todo, handle specific errors
		http.Error(w, "Could not proxy", http.StatusInternalServerError)
	}

	writeResponse(w, resp)
}

func writeResponse(w http.ResponseWriter, resp *http.Response) {
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	io.Copy(w, resp.Body)
	resp.Body.Close()
}
