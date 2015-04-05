package main

import (
	"log"
	"net/http"
	"runtime"
)

func main() {
	// Use all cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	config := new(config)
	config.load()

	s3 := &s3{conf: config}

	proxy := newProxy(config, s3)

	// Simply ok favicon requests
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/", proxy.handler)

	log.Printf("Front v%v listening on port %v ...", VERSION, config.Port)
	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Panic(err)
	}
}
