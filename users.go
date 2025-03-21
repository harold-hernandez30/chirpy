package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/harold-hernandez30/chirpy/internal/auth"
	"github.com/harold-hernandez30/chirpy/internal/database"
)

func (cfg *apiConfig) handleUserLogin(res http.ResponseWriter, req *http.Request) {
	type userCredentialsParam struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	
	decoder := json.NewDecoder(req.Body)
	requestParam := userCredentialsParam{}
	decodeError := decoder.Decode(&requestParam)

	if decodeError != nil {
		log.Printf("unable to decode request body: %s", req.Body)
		res.WriteHeader(500)
		return
	}

	user, findUserErr := cfg.db.FindUser(req.Context(), requestParam.Email)

	if findUserErr != nil {

		log.Printf("user with email '%s' not found", requestParam.Email)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	hashPasswordErr := auth.CheckPasswordHash(requestParam.Password, user.HashedPassword)

	if hashPasswordErr != nil {
		
		log.Printf("password did not match for %s", requestParam.Email)
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("Incorrect email or password"))
		return
	}

	expireTimeInSeconds := 1 * 60 * 60 //1hr

	if requestParam.ExpiresInSeconds >= 0 {
		expireTimeInSeconds = requestParam.ExpiresInSeconds
	}

	jwtRes, makeJwtErr := auth.MakeJWT(user.ID, cfg.secret, time.Duration(expireTimeInSeconds * time.Now().Second()))

	if makeJwtErr != nil {
		
		log.Printf("Something went wrong: %s", makeJwtErr)
		res.WriteHeader(500)
		return
	}

	taggedUser := MapToTaggedUser(user)

	taggedUser.Token = jwtRes
	newUserBytes, marshalErr := json.Marshal(taggedUser)

	if marshalErr != nil {
		log.Printf("Something went wrong: %s", marshalErr)
		res.WriteHeader(500)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(newUserBytes)

}

func (cfg *apiConfig) handleDeleteAllUsers(res http.ResponseWriter, req *http.Request) {

	err := cfg.db.DeleteAllUsers(req.Context())

	if err != nil {
		log.Printf("unable to delete all users: %s", err)
	}

	res.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) handleUserCreate(res http.ResponseWriter, req *http.Request) {
	type userCreateParams struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	requestParam := userCreateParams{}
	decodeError := decoder.Decode(&requestParam)

	if decodeError != nil {
		log.Printf("unable to decode request body: %s", req.Body)
		res.WriteHeader(500)
		return
	}

	hashedPassword, hashPasswordErr := auth.HashPassword(requestParam.Password)

	if hashPasswordErr != nil {

		log.Printf("unable to hash password.\n")
		res.WriteHeader(500)
		return
	}

	userParams := database.CreateUserParams {
		Email: requestParam.Email,
		HashedPassword: hashedPassword,
	}

	newUser, createUserErr := cfg.db.CreateUser(req.Context(), userParams)

	if createUserErr != nil {
		
		log.Printf("unable to save user in the database: %s", createUserErr)
		res.WriteHeader(500)
		return
	}

	taggedUser := MapToTaggedUser(newUser)
	newUserBytes, marshalErr := json.Marshal(taggedUser)

	if marshalErr != nil {
		log.Printf("Something went wrong: %s", marshalErr)
		res.WriteHeader(500)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(201)
	res.Write(newUserBytes)

}