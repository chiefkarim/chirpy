package main

import (
	"log"
	"net/http"
	"os"
)

func (config *apiConfig) reset(response http.ResponseWriter, request *http.Request) {
	config.fileServerHits.Store(0)
	if os.Getenv("PLATFORM") != "dev" {
		response.WriteHeader(401)
		return
	}
	err := config.dbQueries.DeleteAllUsers(request.Context())
	if err != nil {
		log.Printf("Error deleting all users:%v\n", err)
	}
	response.WriteHeader(200)
}
