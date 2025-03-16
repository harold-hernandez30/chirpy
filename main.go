package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func (res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}


func main() {
	mux := http.NewServeMux()
	
	sessionConfig := &apiConfig {
		fileserverHits: atomic.Int32{},
	}
	handler := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	fileServerHandler := http.FileServer(http.Dir("."))
	strippedPrefixHandler := http.StripPrefix("/app", fileServerHandler)

	mux.HandleFunc("GET /admin/metrics", sessionConfig.handlePrintMetrics)
	mux.HandleFunc("POST /admin/reset", sessionConfig.handleResetMetrics)
	
	mux.Handle("/app/", sessionConfig.middlewareMetricsInc(strippedPrefixHandler))

	mux.HandleFunc("GET /api/healthz", handlerHealth)
	
	listenAndServeErr := handler.ListenAndServe()

	if listenAndServeErr != nil {
		log.Fatalf("could not start server: %v", listenAndServeErr)
	}
}

func (cfg *apiConfig) handlePrintMetrics(res http.ResponseWriter, req *http.Request) {
	header := res.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	content := fmt.Sprintf(`
	<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	res.Write([]byte(content))
}

func (cfg *apiConfig) handleResetMetrics(res http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = atomic.Int32{}
}

func handlerHealth(res http.ResponseWriter, req *http.Request) {
	header := res.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	content := []byte(http.StatusText(http.StatusOK))
	res.Write(content)

}