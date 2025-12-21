package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	serverMux := http.NewServeMux()
	server := &http.Server{Addr: ":8080", Handler: serverMux}

	serverMux.Handle("/app/", logger(http.StripPrefix("/app/", http.FileServer(http.Dir("./http")))))
	serverMux.HandleFunc("/healthz", func(response http.ResponseWriter, req *http.Request) {
		response.Header().Add("Content-Type", "text/plain; charset=utf-8")
		response.WriteHeader(200)
		response.Write([]byte("OK"))
	})
	log.Fatal(server.ListenAndServe())
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method + ":" + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
