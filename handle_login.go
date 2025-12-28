package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/chiefkarim/chirpy/internal/auth"
)

func (config *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	var user User
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

	respondWithJSON(w, http.StatusOK, UserDetails{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
		Email:     dbUser.Email,
	})
}
