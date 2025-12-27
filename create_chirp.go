package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/chiefkarim/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chrip struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (config *apiConfig) createChirp(response http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body    string `json:"body"`
		User_id string `json:"user_id"`
	}

	decoder := json.NewDecoder(request.Body)
	var chirp parameters
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding request body %v", err)
		message, err := json.Marshal(map[string]string{"error": "Something went wrong"})
		if err != nil {
			log.Printf("Error Marshaling response error message %v", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.WriteHeader(http.StatusBadRequest)
		response.Write(message)
		return
	}

	if len(chirp.Body) > 140 {
		message, err := json.Marshal(map[string]string{"error": "Chirp is too long"})
		if err != nil {
			log.Printf("Error Marshaling response error message %v", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.WriteHeader(http.StatusBadRequest)
		response.Write(message)
		return
	}

	reg, err := regexp.Compile("(?i)kerfuffle|sharbert|fornax")
	if err != nil {
		log.Printf("Error Marshaling response error message %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	cleaned := reg.ReplaceAll([]byte(chirp.Body), []byte("****"))

	userid, err := uuid.Parse(chirp.User_id)
	if err != nil {
		message, err := json.Marshal(map[string]string{"error": "please enter a valide user id"})
		if err != nil {
			log.Printf("Error Marshaling response error message %v", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.WriteHeader(http.StatusBadRequest)
		response.Write(message)
		return
	}

	createdChirp, err := config.dbQueries.CreateChirp(request.Context(), database.CreateChirpParams{Body: string(cleaned), UserID: userid})
	chripResponse, err := json.Marshal(Chrip{
		Id:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt.Time,
		UpdatedAt: createdChirp.UpdatedAt.Time,
		UserId:    createdChirp.UserID,
		Body:      createdChirp.Body,
	})
	if err != nil {
		log.Printf("Error Marshaling response error message %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
	response.Write(chripResponse)
}
