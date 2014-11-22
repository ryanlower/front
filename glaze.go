package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", rootHandler)

	fmt.Println("Listening to glaze on port " + port + "...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Panic(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fmt.Fprintf(w, "glazed")
}
