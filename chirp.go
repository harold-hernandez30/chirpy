package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/harold-hernandez30/chirpy/internal/database"
)

const CHIRP_MAX_LENGTH = 140

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     string    `json:"body"`
	UserID string `json:"user_id"`
}

func MapToTaggedChirp(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserID: string(dbChirp.UserID.String()),
	}
}

type ChirpCreateParams struct {
	Body string `json:"body"`
}

func decodeChirp(req *http.Request) (*ChirpCreateParams, error) {

	params := ChirpCreateParams{}
	decoder := json.NewDecoder(req.Body)
	decodeErr := decoder.Decode(&params)

	if decodeErr != nil {
		return nil, decodeErr
	}

	return &params, nil
}

func validateChirpCreate(res http.ResponseWriter, message string) error {
	if len(message) > CHIRP_MAX_LENGTH  {
		type ErrorResponse struct {
			Error string `json:"error"`
		}
		res.WriteHeader(400)
		errorMessage := ErrorResponse{
			Error: "Chirp is too long",
		}

		errorBytes, err := json.Marshal(errorMessage)
		res.Write(errorBytes)
		return err
	}

	return nil
}