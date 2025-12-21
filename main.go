package main

import (
	"log"
	"net/http"
)

func main() {
	serverMux := http.NewServeMux()
	server := &http.Server{Addr: ":8080", Handler: serverMux}
	serverMux.Handle("/", http.FileServer(http.Dir("./http")))
	log.Fatal(server.ListenAndServe())
}
