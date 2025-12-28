package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
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

	fileSystemHandler := config.middlewareMetricInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./http"))))
	serverMux.Handle("/app/", fileSystemHandler)

	serverMux.HandleFunc("POST /api/users", config.createUser)
	serverMux.HandleFunc("POST /api/login", config.loginUser)
	serverMux.HandleFunc("GET /api/healthz", config.healthz)
	serverMux.HandleFunc("POST /api/chirps", config.createChirp)
	serverMux.HandleFunc("GET /api/chirps", config.getAllChirps)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", config.getChirp)

	serverMux.HandleFunc("GET /admin/metrics", config.metric)
	serverMux.HandleFunc("POST /admin/reset", config.reset)

	server.Handler = logger(serverMux)
	log.Fatal(server.ListenAndServe())
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method + ":" + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
