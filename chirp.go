package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/harold-hernandez30/chirpy/internal/database"
)

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