package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (config *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	serverMux := http.NewServeMux()
	server := &http.Server{Addr: ":8080", Handler: serverMux}
	config := apiConfig{}

	serverMux.Handle("/app/", config.middlewareMetricInc(logger(http.StripPrefix("/app/", http.FileServer(http.Dir("./http"))))))
	serverMux.HandleFunc("GET /api/healthz", func(response http.ResponseWriter, request *http.Request) {
		response.Header().Add("Content-Type", "text/plain; charset=utf-8")
		response.WriteHeader(200)
		response.Write([]byte("OK"))
	})
	serverMux.HandleFunc("POST /api/validate_chirp", func(response http.ResponseWriter, request *http.Request) {
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

		message, err := json.Marshal(map[string]bool{"valid": true})
		if err != nil {
			log.Printf("Error Marshaling response error message %v", err)
			response.WriteHeader(500)
			return
		}
		response.WriteHeader(200)
		response.Write(message)
	})
	serverMux.HandleFunc("GET /admin/metrics", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(200)
		response.Header().Add("Content-Type", "text/html")
		fmt.Fprintf(response, `<html><body><h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, config.fileServerHits.Load())
	})
	serverMux.HandleFunc("POST /admin/reset", func(response http.ResponseWriter, request *http.Request) {
		config.fileServerHits.Store(0)
		response.WriteHeader(200)
	})
	log.Fatal(server.ListenAndServe())
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method + ":" + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
