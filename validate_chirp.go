package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

func (config *apiConfig) validateChirp(response http.ResponseWriter, request *http.Request) {
	type Chirp struct {
		Body string
	}

	decoder := json.NewDecoder(request.Body)
	var chirp Chirp
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding request body %v", err)
		message, err := json.Marshal(map[string]string{"error": "Something went wrong"})
		if err != nil {
			log.Printf("Error Marshaling response error message %v", err)
			response.WriteHeader(500)
			return
		}
		response.WriteHeader(400)
		response.Write(message)
		return
	}

	if len(chirp.Body) > 140 {
		message, err := json.Marshal(map[string]string{"error": "Chirp is too long"})
		if err != nil {
			log.Printf("Error Marshaling response error message %v", err)
			response.WriteHeader(500)
			return
		}
		response.WriteHeader(400)
		response.Write(message)
		return
	}
	reg, err := regexp.Compile("(?i)kerfuffle|sharbert|fornax")
	if err != nil {
		log.Printf("Error Marshaling response error message %v", err)
		response.WriteHeader(500)
		return
	}
	cleaned := reg.ReplaceAll([]byte(chirp.Body), []byte("****"))

	message, err := json.Marshal(map[string]string{"cleaned_body": string(cleaned)})
	if err != nil {
		log.Printf("Error Marshaling response error message %v", err)
		response.WriteHeader(500)
		return
	}
	response.WriteHeader(200)
	response.Write(message)
}
