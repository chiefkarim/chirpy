package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// define params
type upgradeUserParams struct {
	Event string `json:"event"`
	Data  struct {
		UserId string `json:"user_id"`
	} `json:"data"`
}

func (config *apiConfig) UpgradeUser(response http.ResponseWriter, request *http.Request) {
	// read payload return bad request if can't read
	decoder := json.NewDecoder(request.Body)
	var payload upgradeUserParams
	err := decoder.Decode(&payload)
	if err != nil {
		respondWithJSON(response, http.StatusBadRequest, map[string]string{"error": "request payload should have data with user_id and event type"})
		return
	}

	// if the event is type user.upgraded  if not return 204
	if payload.Event != "user.upgraded" {
		respondWithJSON(response, http.StatusNoContent, nil)
		return
	}

	// return bad request if it doesn't match params
	if payload.Data.UserId == "" {
		respondWithJSON(response, http.StatusBadRequest, map[string]string{"error": "request payload should have data with user_id and event type"})
		return
	}

	// check user id if it's valid otherwise return bad request
	userId, err := uuid.Parse(payload.Data.UserId)
	if err != nil {
		respondWithJSON(response, http.StatusBadRequest, map[string]string{"error": "request payload should have data with user_id and event type"})
		return
	}
	// upgrade user in the database. if the returned is err return 404 otherwise 204
	_, err = config.dbQueries.UpgradeUser(request.Context(), userId)
	if err != nil {
		respondWithJSON(response, http.StatusNotFound, nil)
		return
	}
	respondWithJSON(response, http.StatusNoContent, nil)
}
