package main

import (
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
	serverMux.HandleFunc("/healthz", func(response http.ResponseWriter, req *http.Request) {
		response.Header().Add("Content-Type", "text/plain; charset=utf-8")
		response.WriteHeader(200)
		response.Write([]byte("OK"))
	})
	serverMux.HandleFunc("/metrics", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(200)
		fmt.Fprintf(response, "Hits: %d", config.fileServerHits.Load())
	})
	serverMux.HandleFunc("/reset", func(response http.ResponseWriter, request *http.Request) {
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
