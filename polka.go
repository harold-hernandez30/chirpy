package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlePolkaWebHookEvents(res http.ResponseWriter, req *http.Request) {

	type EventData struct {
		UserID string `json:"user_id"`
	}

	type WebhookEvent struct {
		Event string `json:"event"`
		Data EventData `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	event := WebhookEvent{}
	decodeErr := decoder.Decode(&event)

	if decodeErr != nil {
		fmt.Printf("could not decode web hook event: %s\n", decodeErr)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Event != "user.upgraded" {

		fmt.Printf("unsupported event: %s\n", event.Event)
		res.WriteHeader(http.StatusNoContent)
		return
	}

	// user.upgraded handling

	userID, uuidParseErr := uuid.Parse(event.Data.UserID)

	if uuidParseErr != nil {

		fmt.Printf("could not parse user_id: %s\n", uuidParseErr)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	user, getUserErr := cfg.db.GetUser(req.Context(), userID)

	if getUserErr != nil {
		fmt.Printf("could not find user: %s\n", getUserErr)
		res.WriteHeader(http.StatusNotFound)
		return
	}


	upgradeUserToChirpyRedErr := cfg.db.UpgradeUserToChirpyRed(req.Context(), user.ID)

	if upgradeUserToChirpyRedErr != nil {
		fmt.Printf("could not upgrade user to Chirpy Red: %s\n", upgradeUserToChirpyRedErr)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)

}