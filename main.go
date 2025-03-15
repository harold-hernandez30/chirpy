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

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mux.HandleFunc("/healthz", handlerHealth)

	listenAndServeErr := handler.ListenAndServe()

	if listenAndServeErr != nil {
		log.Fatalf("could not start server: %v", listenAndServeErr)
	}
}

func handlerHealth(res http.ResponseWriter, req *http.Request) {
	header := res.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	content := []byte(http.StatusText(http.StatusOK))
	res.Write(content)

}