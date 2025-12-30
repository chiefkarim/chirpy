package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (config *apiConfig) getAllChirps(response http.ResponseWriter, request *http.Request) {
	urlParams := request.URL.Query()
	authorId := urlParams.Get("author_id")

	if authorId != "" {
		parsedAuthorId, err := uuid.Parse(authorId)
		if err != nil {
			respondWithJSON(response, http.StatusUnauthorized, map[string]string{"error": "Please provide valid user id query"})
			return
		}
		rows, err := config.dbQueries.GetChirpsByUserId(request.Context(), parsedAuthorId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithJSON(response, http.StatusNotFound, map[string]string{"error": "No chirps found for the given author id"})
				return
			}
			respondWithJSON(response, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			return
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
		respondWithJSON(response, http.StatusOK, chirps)
		return
	}

	rows, err := config.dbQueries.GetAllChirps(request.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(response, http.StatusNotFound, map[string]string{"error": "No chirps found for the given author id"})
			return
		}
		respondWithJSON(response, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
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
	respondWithJSON(response, http.StatusOK, chirps)
}
