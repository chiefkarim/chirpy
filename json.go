package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(response http.ResponseWriter, status int, payload any) {
	response.Header().Add("Content-type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error masrshaling JSON:%v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(status)
	response.Write(data)
}

func JSONResponse5OO(w http.ResponseWriter) {
	message := map[string]string{"error": "internal server error"}
	respondWithJSON(w, http.StatusInternalServerError, message)
}
