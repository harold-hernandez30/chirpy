package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	handler := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))

	listenAndServeErr := handler.ListenAndServe()

	if listenAndServeErr != nil {
		log.Fatalf("could not start server: %v", listenAndServeErr)
	}
}