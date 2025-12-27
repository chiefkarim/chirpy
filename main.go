package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/chiefkarim/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries      *database.Queries
}

func (config *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("DB_URL environment variable not found!")
	}

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatalf("Somthing went wrong connecting to Databse %s: \n%v", db_url, err)
	}
	dbQueries := database.New(db)
	config := apiConfig{dbQueries: dbQueries}

	serverMux := http.NewServeMux()
	server := &http.Server{Addr: ":8080", Handler: serverMux}
	fmt.Printf("server listening on: http://localhost:%s/app/\n", strings.ReplaceAll(server.Addr, ":", ""))

	fileSystemHandler := config.middlewareMetricInc(logger(http.StripPrefix("/app/", http.FileServer(http.Dir("./http")))))
	serverMux.Handle("/app/", fileSystemHandler)

	serverMux.HandleFunc("POST /api/users", config.createUser)
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
		if os.Getenv("PLATFORM") != "dev" {
			response.WriteHeader(401)
			return
		}
		err := config.dbQueries.DeleteAllUsers(request.Context())
		if err != nil {
			log.Printf("Error deleting all users:%v\n", err)
		}
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
