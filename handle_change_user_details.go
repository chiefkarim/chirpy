package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/chiefkarim/chirpy/internal/auth"
	"github.com/chiefkarim/chirpy/internal/database"
)

type Params struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (config *apiConfig) ChangeUserDetails(response http.ResponseWriter, request *http.Request) {
	bearerToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithJSON(response, http.StatusUnauthorized, map[string]string{"error": "Plase provide valid accesstoken header"})
		return
	}

	var requestPayload Params
	body, err := io.ReadAll(request.Body)
	if err != nil {
		respondWithJSON(response, http.StatusBadRequest, map[string]string{"error": "Plase provide valid email and password"})
		return
	}

	err = json.Unmarshal(body, &requestPayload)
	if err != nil {
		respondWithJSON(response, http.StatusBadRequest, map[string]string{"error": "Plase provide valid email and password"})
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, config.hashKey)
	if err != nil {
		respondWithJSON(response, http.StatusUnauthorized, map[string]string{"error": "Plase provide valid accesstoken header"})
		return
	}

	hashedPassword, err := auth.HashPassword(requestPayload.Password)
	if err != nil {
		JSONResponse5OO(response)
		return
	}
	newUser, err := config.dbQueries.ChangeUserDetails(request.Context(), database.ChangeUserDetailsParams{
		Email:          requestPayload.Email,
		ID:             userId,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		JSONResponse5OO(response)
		return
	}

	respondWithJSON(response, http.StatusOK, UserDetails{
		ID:        newUser.ID,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt.Time,
		UpdatedAt: newUser.UpdatedAt.Time,
	})
}
