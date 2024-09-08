package main

import (
	"log"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/signatures"
)

const (
	ListenAddress = ":8080"
	// TODO: add further configuration parameters here ...
)

func main() {
	server := http.NewServeMux()

	storage := signatures.NewMemory()

	server.Handle("/signatures/", http.StripPrefix("/signatures", signatures.NewHandler(server, storage)))

	log.Fatal(http.ListenAndServe(ListenAddress, server))
}
