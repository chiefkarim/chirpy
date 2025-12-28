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

type User struct {
	Email    string `json:"email"`
	Paswword string `json:"password"`
}

type UserDetails struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (config *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("500:Error reading body for createUser %v\n", err)
		JSONResponse5OO(w)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("500:Error unmarshaling body for createUser %v\n", err)
		JSONResponse5OO(w)
		return
	}

	if user.Email == "" || user.Paswword == "" {
		log.Println("401:Error request doesn't have email or password fields!")
		message := map[string]string{"error": "must provide email and password"}
		respondWithJSON(w, http.StatusBadRequest, message)
		return
	}

	hashedPassword, err := auth.HashPassword(user.Paswword)
	if err != nil {
		log.Printf("500:Error hashing password %v\n", err)
		JSONResponse5OO(w)
	}

	createdUser, err := config.dbQueries.CreateUser(r.Context(), database.CreateUserParams{Email: user.Email, HashedPassword: hashedPassword})
	if err != nil {
		log.Printf("Error creating user with email %s\n%v", user.Email, err)
		JSONResponse5OO(w)
		return
	}

	respondWithJSON(w, http.StatusCreated, UserDetails{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt.Time,
		UpdatedAt: createdUser.UpdatedAt.Time,
		Email:     createdUser.Email,
	})
}
