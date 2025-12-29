package main

import (
	"log"
	"net/http"
	"time"

	"github.com/chiefkarim/chirpy/internal/auth"
)

func (config *apiConfig) RefreshToken(response http.ResponseWriter, request *http.Request) {
	refreshToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		log.Printf("500:Error reading body for createUser %v\n", err)
		JSONResponse5OO(response)
		return
	}

	DBRefreshToken, err := config.dbQueries.GetRefreshToken(request.Context(), refreshToken)
	if err != nil {
		log.Printf("500:Error reading body for createUser %v\n", err)
		JSONResponse5OO(response)
		return
	}

	if DBRefreshToken.ExpiresAt.Before(time.Now().UTC()) || DBRefreshToken.RevokedAt.Time.Before(time.Now().UTC()) {
		respondWithJSON(response, http.StatusUnauthorized, map[string]string{"error": "Your refresh token has expired"})
		return
	}

	accessToken, err := auth.MakeJWT(DBRefreshToken.UserID, config.hashKey, 1*time.Hour)
	if err != nil {
		JSONResponse5OO(response)
		return
	}
	respondWithJSON(response, http.StatusOK, map[string]string{"token": accessToken})
}
