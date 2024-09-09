package main

import (
	"log"
	"net/http"

	_ "embed"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/docs"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/signatures"
)

const (
	ListenAddress = ":8080"
)

func main() {
	server := http.NewServeMux()

	storage := signatures.NewMemory()

	server.Handle("/", docs.NewHandler())
	// server.Handle("/docs/", http.StripPrefix("/docs", docs.NewHandler()))
	server.Handle("/signatures/", http.StripPrefix("/signatures", signatures.NewHandler(storage)))

	log.Fatal(http.ListenAndServe(ListenAddress, server))
}
