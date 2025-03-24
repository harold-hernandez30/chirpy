package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/harold-hernandez30/chirpy/internal/auth"
	"github.com/harold-hernandez30/chirpy/internal/database"
)

func (cfg *apiConfig) handleRefresh(res http.ResponseWriter, req *http.Request) {
	token, bearerTokenErr := auth.GetBearerToken(req.Header)

	if bearerTokenErr != nil {
		fmt.Printf("unable to parse refresh token: %s\n", bearerTokenErr)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	userAndTokenRow, getUserFromRefreshTokenErr := cfg.db.GetUserFromRefreshToken(req.Context(), token)
	if getUserFromRefreshTokenErr != nil {
		
		fmt.Printf("user not found: %s\n", getUserFromRefreshTokenErr)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("user and row found")

	if time.Now().After(userAndTokenRow.RefreshToken.ExpiresAt) {
		
		fmt.Printf("user refresh token expired: %s\n", userAndTokenRow.RefreshToken.RevokedAt.Time)
		res.WriteHeader(http.StatusUnauthorized)
	}

	
	fmt.Printf("refresh token is valid")

	type tokenResponse struct {
		Token string `json:"token"`
	}

	newRefreshToken, _ := auth.MakeRefreshToken()


	refreshTokenParams := database.CreateRefreshTokenParams {
		Token: newRefreshToken,
		UserID: userAndTokenRow.User.ID,
	}
	
	
	_, createRefresTokenErr := cfg.db.CreateRefreshToken(req.Context(), refreshTokenParams)

	if createRefresTokenErr != nil {
		log.Printf("Something went wrong: %s", createRefresTokenErr)
		res.WriteHeader(500)
		return
	}

	newAccessToken, makeJWTErr := auth.MakeJWT(userAndTokenRow.RefreshToken.UserID, cfg.secret, 1 * time.Hour)


	if makeJWTErr != nil {
		
		log.Printf("unable to create jwt: %s", makeJWTErr)
		res.WriteHeader(500)
		return
	}


	resBody := tokenResponse {
		Token: newAccessToken,
	}
	if resBodyInBytes, marshalErr := json.Marshal(&resBody); marshalErr == nil {

		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(resBodyInBytes)
	} else {
		
		fmt.Printf("something went wrong: %s\n", marshalErr)
		res.WriteHeader(http.StatusInternalServerError)
	}

}

func (cfg *apiConfig) handleUserLogin(res http.ResponseWriter, req *http.Request) {
	type userCredentialsParam struct {
		Email string `json:"email"`
		Password string `json:"password"`
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


	jwtRes, makeJwtErr := auth.MakeJWT(user.ID, cfg.secret, 1 * time.Hour)

	if makeJwtErr != nil {
		
		log.Printf("Something went wrong: %s", makeJwtErr)
		res.WriteHeader(500)
		return
	}

	taggedUser := MapToTaggedUser(user)
	refreshToken, makeRefreshTokenErr := auth.MakeRefreshToken()

	if makeRefreshTokenErr != nil {
		log.Printf("Something went wrong: %s", makeRefreshTokenErr)
		res.WriteHeader(500)
		return
	}


	taggedUser.AccessToken = jwtRes
	
	refreshTokenParams := database.CreateRefreshTokenParams {
		Token: refreshToken,
		UserID: user.ID,
	}
	
	refreshTokenRow, createRefresTokenErr := cfg.db.CreateRefreshToken(req.Context(), refreshTokenParams)

	if createRefresTokenErr != nil {
		log.Printf("Something went wrong: %s", createRefresTokenErr)
		res.WriteHeader(500)
		return
	}

	taggedUser.RefreshToken = refreshTokenRow.Token

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