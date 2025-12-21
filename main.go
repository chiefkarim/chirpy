package main

import (
	"log"
	"net/http"
)

func main() {
	serverMux := http.NewServeMux()
	server := &http.Server{Addr: ":8080", Handler: serverMux}

	serverMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./http"))))
	serverMux.HandleFunc("/healthz", func(response http.ResponseWriter, req *http.Request) {
		response.Header().Add("Content-Type", "text/plain; charset=utf-8")
		response.WriteHeader(200)
		response.Write([]byte("OK"))
	})

	log.Fatal(server.ListenAndServe())
}
