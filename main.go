package main

import (
	"log"
	"net/http"
)

func main() {
	addr := "localhost:8080"
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.NotFound(w, req)
		return
	})
	log.Printf("Server listening on %s\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
