package main

import (
	"encoding/json"
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
	mux.HandleFunc("POST /api/validate_chirp", handleChirpValidate)
	
	mux.Handle("/app/", sessionConfig.middlewareMetricsInc(strippedPrefixHandler))

	mux.HandleFunc("GET /api/healthz", handlerHealth)
	
	listenAndServeErr := handler.ListenAndServe()

	if listenAndServeErr != nil {
		log.Fatalf("could not start server: %v", listenAndServeErr)
	}
}

func handleChirpValidate(res http.ResponseWriter, req *http.Request) {
	type Params struct {
		Body string `json:"body"`
	}

	type ErrorResponse struct {
		Error string `json:"error"`
	}


	decoder := json.NewDecoder(req.Body)
	reqParams := Params{}

	err := decoder.Decode(&reqParams)

	res.Header().Set("Content-Type", "application/json")

	if err != nil {

		log.Printf("error decoding parameters: %s", err)
		res.WriteHeader(500)
		errorMessage := ErrorResponse{
			Error: "Something went wrong",
		}

		errorBytes, _ := json.Marshal(errorMessage)
		res.Write(errorBytes)
		return
	}

	if len(reqParams.Body) > 140 {
		res.WriteHeader(400)
		errorMessage := ErrorResponse{
			Error: "Chirp is too long",
		}

		errorBytes, _ := json.Marshal(errorMessage)
		res.Write(errorBytes)
		return
	}


	res.WriteHeader(200)


	type ValidChirp struct {
		Valid bool `json:"valid"`
	}


	validChirpResponse, _ := json.Marshal(ValidChirp{Valid: true})

	res.Write(validChirpResponse)
}

func handlerHealth(res http.ResponseWriter, req *http.Request) {
	header := res.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	content := []byte(http.StatusText(http.StatusOK))
	res.Write(content)

}