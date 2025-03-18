package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/harold-hernandez30/chirpy/internal/database"
)

func (cfg *apiConfig) handleChirpCreate(res http.ResponseWriter, req *http.Request) {
	type params struct {
		Body string `json:"body"`
		UserId string `json:"user_id"`
	}

	reqParams := params{}

	decoder := json.NewDecoder(req.Body)
	decodeErr := decoder.Decode(&reqParams)

	if decodeErr != nil {
		log.Printf("unable to decode: %s", decodeErr)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	nullUUID := uuid.NullUUID {}
	uuidErr := nullUUID.Scan(reqParams.UserId)

	if uuidErr != nil {
		log.Printf("invalid uuid: %s", uuidErr)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	user, getUserErr := cfg.db.GetUser(req.Context(), nullUUID.UUID)

	if getUserErr != nil {
		
		log.Printf("error getting user (UUID: %s): %s", nullUUID.UUID.String(), getUserErr)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	type ErrorResponse struct {
		Error string `json:"error"`
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

	chirpParams := database.CreateChirpParams {
		Body: reqParams.Body,
		UserID: user.ID,
	}
	newChirp, createChirpErr := cfg.db.CreateChirp(req.Context(), chirpParams)

	if createChirpErr != nil {

		log.Printf("error saving chirp: %s", createChirpErr)
		res.WriteHeader(http.StatusBadRequest)
	}

	taggedChirp := MapToTaggedChirp(newChirp)

	newChirpBytes, marchalErr := json.Marshal(taggedChirp)

	
	if marchalErr != nil {
		log.Printf("Something went wrong: %s", marchalErr)
		res.WriteHeader(500)
		return
	}
	
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(201)
	res.Write(newChirpBytes)

}