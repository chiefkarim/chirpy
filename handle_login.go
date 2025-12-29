package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/chiefkarim/chirpy/internal/auth"
	"github.com/chiefkarim/chirpy/internal/database"
	"github.com/google/uuid"
)

type params struct {
	Email    string `json:"email"`
	Paswword string `json:"password"`
}
type LoginUserDetails struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (config *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	var user params
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("500:Error reading body for loginUser %v\n", err)
		JSONResponse5OO(w)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("500:Error unmarshaling body for loginUser %v\n", err)
		JSONResponse5OO(w)
		return
	}

	dbUser, err := config.dbQueries.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		log.Printf("401:Error getting user by email %v\n", err)
		message := map[string]string{"error": "Incorrect email or password"}
		respondWithJSON(w, http.StatusUnauthorized, message)
		return
	}

	isValid, err := auth.CheckPasswordHash(user.Paswword, dbUser.HashedPassword)
	if err != nil {
		log.Printf("500:Error checking password hash %v\n", err)
		JSONResponse5OO(w)
		return
	}

	if isValid != true {
		log.Printf("401:Error checking password hash %v\n", err)
		message := map[string]string{"error": "Incorrect email or password"}
		respondWithJSON(w, http.StatusUnauthorized, message)
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, config.hashKey, 1*time.Hour)
	if err != nil {
		JSONResponse5OO(w)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	log.Printf("refresh token, %s", refreshToken)
	if err != nil {
		JSONResponse5OO(w)
		return
	}
	_, err = config.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
		UserID:    dbUser.ID,
	})
	if err != nil {
		JSONResponse5OO(w)
		return
	}
	respondWithJSON(w, http.StatusOK, LoginUserDetails{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}
