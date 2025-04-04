package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/harold-hernandez30/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	AccessToken	  string	`json:"token"`
	RefreshToken string `json:"refresh_token"`
	UserRedMember bool `json:"is_chirpy_red"`
}

func MapToTaggedUser(dbUser database.User) User {
	return User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
		UserRedMember: dbUser.IsChirpyRed.Bool,
	}
}