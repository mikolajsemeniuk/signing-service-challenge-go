package main

import (
	"log"
	"net/http"
	"time"

	_ "embed"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/docs"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/signature"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Listen       string        `envconfig:"LISTEN"        required:"true"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT"  default:"5s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"5s"`
	IdleTimeout  time.Duration `envconfig:"IDLE_TIMEOUT"  default:"30s"`
}

func main() {
	var config config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}

	storage := signature.NewMemory()

	router := http.NewServeMux()

	router.Handle("/", docs.NewHandler())
	router.Handle("/signature/", http.StripPrefix("/signature", signature.NewHandler(storage)))

	server := &http.Server{
		Addr:         config.Listen,
		Handler:      router,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	log.Printf("Server starting on %s", config.Listen)
	log.Fatal(server.ListenAndServe())
}
