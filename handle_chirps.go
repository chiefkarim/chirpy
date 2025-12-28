package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (config *apiConfig) getAllChirps(response http.ResponseWriter, request *http.Request) {
	rows, err := config.dbQueries.GetAllChirps(request.Context())
	if err != nil {
		message, err := json.Marshal(map[string]string{"error": "error getting all chirps."})
		if err != nil {
			log.Printf("Error Marshaling response error message %v", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Error getting all chirps from db. %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(message)
	}
	chirps := make([]Chirp, 0, len(rows))
	for _, row := range rows {
		chirps = append(chirps, Chirp{
			Id:        row.ID,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			UserId:    row.UserID,
			Body:      row.Body,
		})
	}
	parsedChirps, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("Error Marshaling response error message %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(parsedChirps)
}
