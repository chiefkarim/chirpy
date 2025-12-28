package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (config *apiConfig) getChirp(response http.ResponseWriter, request *http.Request) {
	chirpID := request.PathValue("chirpID")
	parsedChirpId, err := uuid.Parse(chirpID)
	if chirpID == "" || err != nil {
		respondWithJSON(response, http.StatusBadRequest, map[string]string{"error": "must provide valid chirpID in the url!"})
		return
	}

	chirp, err := config.dbQueries.GetChirp(request.Context(), parsedChirpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(response, http.StatusNotFound, map[string]string{"error": "No shirp was found with the given id!"})
			return
		}
		message := map[string]string{"error": "Somthing went wrong!"}
		log.Printf("500: Error,%v\n", err)
		respondWithJSON(response, http.StatusInternalServerError, message)
		return
	}

	parsedChirps, err := json.Marshal(Chirp{
		Id:        chirp.ID,
		UpdatedAt: chirp.UpdatedAt.Time,
		CreatedAt: chirp.CreatedAt.Time,
		UserId:    chirp.UserID,
		Body:      chirp.Body,
	})
	if err != nil {
		message := map[string]string{"error": "Somthing went wrong!"}
		log.Printf("500: Error Marshaling response error message %v", err)
		respondWithJSON(response, http.StatusInternalServerError, message)
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(parsedChirps)
}
