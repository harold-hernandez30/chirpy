package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handleDeleteAllUsers(res http.ResponseWriter, req *http.Request) {

	err := cfg.db.DeleteAllUsers(req.Context())

	if err != nil {
		log.Printf("unable to delete all users: %s", err)
	}

	res.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) handleUserCreate(res http.ResponseWriter, req *http.Request) {
	type reqParams struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	requestParam := reqParams{}
	decodeError := decoder.Decode(&requestParam)

	if decodeError != nil {
		log.Printf("unable to decode request body: %s", req.Body)
		res.WriteHeader(500)
		return
	}



	newUser, createUserErr := cfg.db.CreateUser(req.Context(), requestParam.Email)

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