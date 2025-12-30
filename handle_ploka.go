package main

import (
	"encoding/json"
	"net/http"

	"github.com/chiefkarim/chirpy/internal/auth"
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
	reqApiKey, err := auth.GetAPIKey(request.Header)
	if err != nil || reqApiKey != config.plokaApiKey {
		respondWithJSON(response, http.StatusUnauthorized, map[string]string{"error": "Please provide valid api key in the header"})
		return
	}

	decoder := json.NewDecoder(request.Body)
	var payload upgradeUserParams
	err = decoder.Decode(&payload)
	if err != nil {
		respondWithJSON(response, http.StatusBadRequest, map[string]string{"error": "request payload should have data with user_id and event type"})
		return
	}

	if payload.Event != "user.upgraded" {
		respondWithJSON(response, http.StatusNoContent, nil)
		return
	}

	if payload.Data.UserId == "" {
		respondWithJSON(response, http.StatusBadRequest, map[string]string{"error": "request payload should have data with user_id and event type"})
		return
	}

	userId, err := uuid.Parse(payload.Data.UserId)
	if err != nil {
		respondWithJSON(response, http.StatusBadRequest, map[string]string{"error": "request payload should have data with user_id and event type"})
		return
	}

	_, err = config.dbQueries.UpgradeUser(request.Context(), userId)
	if err != nil {
		respondWithJSON(response, http.StatusNotFound, nil)
		return
	}
	respondWithJSON(response, http.StatusNoContent, nil)
}
