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

func main() {
	config := Config{
		port:                os.Getenv("PORT"),
		allowedContentTypes: "^image/",
	}
	proxy := newProxy(config)

	http.HandleFunc("/", proxy.handler)

	log.Println("Listening to glaze on port " + config.port + "...")
	err := http.ListenAndServe(":"+config.port, nil)
	if err != nil {
		log.Panic(err)
	}
}
