package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type User struct {
	Email string `json:"email"`
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
		fmt.Fprint(w,
			`
		{
			"error":"wrong body format",
			"expected":{
				email:"example@example.com"
			}
		}
		 `)
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

	type UserDetails struct {
		Id          string `json:"id"`
		Created_at  string `json:"created_at"`
		Updateed_at string `json:"updated_at"`
		Email       string `json:"email"`
	}
	returendUser, err := json.Marshal(UserDetails{Id: createdUser.ID.String(), Created_at: createdUser.CreatedAt.Time.String(), Updateed_at: createdUser.UpdatedAt.Time.String(), Email: createdUser.Email})
	if err != nil {
		log.Printf("error unmarshaling body for createUser %v\n", err)
		w.WriteHeader(500)
	}

	w.WriteHeader(201)
	w.Header().Add("Content-type", "application/json")
	fmt.Fprint(w, string(returendUser))
}
