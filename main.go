package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	handler := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	listenAndServeErr := handler.ListenAndServe()

	if listenAndServeErr != nil {
		log.Fatalf("could not start server")
	}
}