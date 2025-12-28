package main

import (
	"fmt"
	"net/http"
)

func (config *apiConfig) metric(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(response, `<html><body><h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, config.fileServerHits.Load())
}
