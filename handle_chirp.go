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
		message, err := json.Marshal(map[string]string{"error": "must provide valid chirpID in the url!"})
		if err != nil {
			log.Printf("500: error while marshaling error message,%v\n", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.WriteHeader(http.StatusBadRequest)
		response.Write(message)
		return
	}
	chirp, err := config.dbQueries.GetChirp(request.Context(), parsedChirpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			message, err := json.Marshal(map[string]string{"error": "No shirp was found with the given id!"})
			if err != nil {
				log.Printf("500: error while marshaling error message,%v\n", err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			response.WriteHeader(http.StatusNotFound)
			response.Write(message)
			return
		}
		message, err := json.Marshal(map[string]string{"error": "Couldn't get chirps for the provided user id!"})
		if err != nil {
			log.Printf("500: error while marshaling error message,%v\n", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(message)
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
		log.Printf("Error Marshaling response error message %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(parsedChirps)
}
