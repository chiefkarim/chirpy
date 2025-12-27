package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Email string `json:"email"`
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
		log.Printf("error reading body for createUser %v\n", err)
		w.WriteHeader(401)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("error unmarshaling body for createUser %v\n", err)
		w.WriteHeader(401)
		fmt.Fprint(w, "Expected a valid JSON\n")
		return
	}

	if user.Email == "" {
		log.Println("Error: request doesn't have email field!")
		w.WriteHeader(401)
		fmt.Fprint(w,
			"{'error':'must provide email'}\n")
		return
	}

	createdUser, err := config.dbQueries.CreateUser(r.Context(), user.Email)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error creating user with email %s\n%v", user.Email, err)
		return
	}

	returendUser, err := json.Marshal(UserDetails{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt.Time,
		UpdatedAt: createdUser.UpdatedAt.Time,
		Email:     createdUser.Email,
	})
	if err != nil {
		log.Printf("error unmarshaling body for createUser %v\n", err)
		w.WriteHeader(500)
	}

	w.WriteHeader(201)
	w.Header().Add("Content-type", "application/json")
	fmt.Fprint(w, string(returendUser))
}
