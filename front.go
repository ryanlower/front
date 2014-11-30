package main

import (
	"log"
	"net/http"
	"os"
)

type Config struct {
	port                string
	allowedContentTypes string // uncompiled regex
}

func envOrDefault(key string, default_value string) string {
	env := os.Getenv(key)
	if env != "" {
		return env
	} else {
		return default_value
	}
}

func main() {
	config := Config{
		port:                os.Getenv("PORT"),
		allowedContentTypes: envOrDefault("ALLOWED_CONTENT_TYPE_REGEX", "^image/"),
	}
	proxy := newProxy(config)

	// Simply ok favicon requests
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/", proxy.handler)

	log.Printf("Front v%v listening on port %v ...", VERSION, config.port)
	err := http.ListenAndServe(":"+config.port, nil)
	if err != nil {
		log.Panic(err)
	}
}
