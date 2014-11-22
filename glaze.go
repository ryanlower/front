package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", proxy)

	log.Println("Listening to glaze on port " + port + "...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Panic(err)
	}
}
