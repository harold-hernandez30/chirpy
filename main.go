package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/harold-hernandez30/chirpy/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)
type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func (res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}

func (cfg *apiConfig) middlewareDevOnly(next func (res http.ResponseWriter, req *http.Request)) func (http.ResponseWriter, *http.Request) {
	return func (res http.ResponseWriter, req *http.Request) {
		if cfg.platform == "dev" {
			next(res, req)
		} else {
			res.WriteHeader(http.StatusForbidden)
			content := []byte(http.StatusText(http.StatusForbidden))
			res.Write(content)
		}
	}
}


func main() {
	// Database setup
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatalf("dbURL must be set")
	}

	db, dbErr := sql.Open("postgres", dbURL)

	if dbErr != nil {
		log.Fatalf("error connecting to database (%s): %s", dbURL, dbErr)
	}


	dbQueries := database.New(db)

	mux := http.NewServeMux()
	
	sessionConfig := &apiConfig {
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		platform: os.Getenv("PLATFORM"),
	}
	handler := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	fileServerHandler := http.FileServer(http.Dir("."))
	strippedPrefixHandler := http.StripPrefix("/app", fileServerHandler)

	mux.HandleFunc("GET /admin/metrics", sessionConfig.handlePrintMetrics)
	mux.HandleFunc("POST /admin/reset", sessionConfig.middlewareDevOnly(sessionConfig.handleDeleteAllUsers))
	mux.HandleFunc("POST /api/validate_chirp", handleChirpValidate)
	mux.HandleFunc("POST /api/users", sessionConfig.handleUserCreate)
	mux.HandleFunc("POST /api/chirps", sessionConfig.handleChirpCreate)
	mux.HandleFunc("GET /api/chirps", sessionConfig.handleGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", sessionConfig.handleGetChirp)
	
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


	res.WriteHeader(http.StatusOK)


	type CleanMessage struct {
		CleanedBody string `json:"cleaned_body"`
	}


	cleanedMessage := cleanMessage(reqParams.Body)
	validChirpResponse, _ := json.Marshal(CleanMessage{CleanedBody: cleanedMessage})

	res.Write(validChirpResponse)
}

func cleanMessage(msg string) string {
	profaneWords := []string{"kerfuffle", "sharbert","fornax"}

	allWords := strings.Split(msg, " ")

	for i, word := range allWords {

		for _, profaneWord := range profaneWords {
			if strings.ToLower(word) == profaneWord {
				allWords[i] = "****"
				continue
			}
			
		}
		
		
	}

	return strings.Join(allWords, " ")
}

func handlerHealth(res http.ResponseWriter, req *http.Request) {
	header := res.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	content := []byte(http.StatusText(http.StatusOK))
	res.Write(content)

}