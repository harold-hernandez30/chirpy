package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/harold-hernandez30/chirpy/internal/database"
)

func (cfg *apiConfig) handleDeleteChirp(res http.ResponseWriter, req *http.Request) {
	if userUUID, validateJWTErr := cfg.getUserFromAuthorization(req); validateJWTErr != nil {
		fmt.Printf("validating JWT error failed: %s\n", validateJWTErr)
		res.WriteHeader(http.StatusUnauthorized)
		content := []byte(http.StatusText(http.StatusUnauthorized))
		res.Write(content)
	} else {
		
		chirpID := req.PathValue("chirpID")
		if len(chirpID) == 0 {
			handleError(
			res,
			fmt.Errorf("chirpID not found"),
			http.StatusBadRequest,
			fmt.Sprintln("chirpID must be provided"))
			return
		}

		uuidChirp, parseChirpIdErr := uuid.Parse(chirpID)

		if parseChirpIdErr != nil {
			fmt.Printf("unable to parse uuid provided: %s\n", parseChirpIdErr)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		
		foundChirp, getChirpErr := cfg.db.GetChirp(req.Context(), uuidChirp)

		if parseChirpIdErr != nil {
			fmt.Printf("chirp does not exist: %s\n", getChirpErr)
			res.WriteHeader(http.StatusNotFound)
			return
		}

		if foundChirp.UserID.String() != userUUID.String() {
			fmt.Printf("chirp not owned by user: %s\n", foundChirp)
			res.WriteHeader(http.StatusForbidden)
			return
		}

		deleteChirpErr := cfg.db.DeleteChirp(req.Context(), foundChirp.ID)

		if deleteChirpErr != nil {
			
			fmt.Printf("unabled to delete chirp: %s\n", deleteChirpErr)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusNoContent)

	}
	
}

func (cfg *apiConfig) handleGetChirp(res http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")

	
	var e error
	if len(chirpID) == 0 {
		handleError(
		res,
		fmt.Errorf("chirpID not found"),
		http.StatusBadRequest,
		fmt.Sprintln("chirpID must be provided"))
		return
	}

	
	nullUUID := uuid.NullUUID {}
	uuidErr := nullUUID.Scan(chirpID)

	if e = handleError(res, 
		uuidErr, 
		http.StatusBadRequest, 
		fmt.Sprintf("invalid uuid: %s", uuidErr)); e != nil {
		return
	}

	foundChirp, getChirpErr := cfg.db.GetChirp(req.Context(), nullUUID.UUID)

	
	if e = handleError(res, 
		getChirpErr, 
		http.StatusNotFound, 
		fmt.Sprintf("chirp with id (%s) not found", chirpID)); e != nil {
		return
	}

	chirpResponse := MapToTaggedChirp(foundChirp)
	chirpResponseByte, marshalErr := json.Marshal(&chirpResponse)

	if marshalErr != nil {
		log.Printf("Something went wrong: %s", marshalErr)
		res.WriteHeader(500)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(chirpResponseByte)

}

func (cfg *apiConfig) handleGetAllChirps(res http.ResponseWriter, req *http.Request) {

	queryParams := req.URL.Query()
	filterAuthorId := queryParams.Get("author_id")
	var e error
	
	if filterAuthorId != "" {
		authorUUID, parseErr := uuid.Parse(filterAuthorId)


		if e = handleError(
			res,
			parseErr,
			http.StatusBadRequest,
			fmt.Sprintf("unable to parse author_id: %s", parseErr)); e != nil {
			return
		}

		chirpsFromUser, chirpsFromUserErr := cfg.db.GetChirpsFromUser(req.Context(), authorUUID)

		
		if e = handleError(
			res,
			chirpsFromUserErr,
			http.StatusBadRequest,
			fmt.Sprintf("unable to get chirps from author_id: %s", chirpsFromUserErr)); e != nil {
			return
		}

		chirpsResponseHelper(chirpsFromUser, res, req)
	} else {

		dbAllChirps, getAllChirpsErr := cfg.db.GetAllChirps(req.Context())

		if e = handleError(
			res,
			getAllChirpsErr,
			http.StatusBadRequest,
			fmt.Sprintf("unable to retrieve chirps from DB: %s", getAllChirpsErr)); e != nil {
			return
		}

		chirpsResponseHelper(dbAllChirps, res, req)
	}
}

func chirpsResponseHelper(chirpsFromDb []database.Chirp, res http.ResponseWriter, req *http.Request) error {

	sortBy := req.URL.Query().Get("sort")

	chirpSlice := []Chirp{}
	for _, dbChirp := range chirpsFromDb {
		taggedChirp := MapToTaggedChirp(dbChirp)
		chirpSlice = append(chirpSlice, taggedChirp)
	}

	if len(sortBy) > 0 {
		sort.Slice(chirpSlice, func(i, j int) bool {
			if (sortBy == "asc") {
				return chirpSlice[j].CreatedAt.After(chirpSlice[i].CreatedAt)
			} else if (sortBy == "desc") {
				return chirpSlice[i].CreatedAt.After(chirpSlice[j].CreatedAt)
			} else {
				return true
			}
		})
	}

	allChirpsBytes, marchalErr := json.Marshal(chirpSlice)

	if marchalErr != nil {
		log.Printf("Something went wrong: %s", marchalErr)
		res.WriteHeader(500)
		return marchalErr
	}
	
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(allChirpsBytes)
	return nil
}

func (cfg *apiConfig) handleChirpCreate(res http.ResponseWriter, req *http.Request) {

	decodedParams, decodeErr := decodeChirp(req)
	var e error
	if e = handleError(res, 
		decodeErr, 
		http.StatusBadRequest, 
		fmt.Sprintf("unable to decode: %s", decodeErr)); e != nil {
		return
	}

	user, getUserErr := cfg.db.GetUser(req.Context(), cfg.currentUserUUID)

	if e = handleError(res, 
		getUserErr, 
		http.StatusBadRequest, 
		fmt.Sprintf("error getting user (UUID: %s): %s", cfg.currentUserUUID.String(), getUserErr)); e != nil {
		return
	}

	
	chirpCreateErr := validateChirpCreate(res, decodedParams.Body)

	if chirpCreateErr != nil {
		return
	}

	chirpParams := database.CreateChirpParams {
		Body: decodedParams.Body,
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
	res.WriteHeader(http.StatusCreated)
	res.Write(newChirpBytes)

}