package main

import (
	"html/template"
	"log"
	"net/http"

	_ "embed"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/signatures"
)

//go:embed openapi.yaml
var docs []byte

//go:embed elements.tpl
var elements string

const (
	ListenAddress = ":8080"
)

func main() {
	server := http.NewServeMux()

	storage := signatures.NewMemory()

	server.HandleFunc("GET /x", func(w http.ResponseWriter, r *http.Request) {
		template.Must(template.New("ui").Parse(elements)).Execute(w, "./docs")
	})
	server.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) { w.Write(docs) })

	server.Handle("/signatures/", http.StripPrefix("/signatures", signatures.NewHandler(server, storage)))

	log.Fatal(http.ListenAndServe(ListenAddress, server))
}
