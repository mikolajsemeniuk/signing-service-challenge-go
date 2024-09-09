package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/docs"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/signature"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Listen       string        `envconfig:"LISTEN" required:"true"`
	ReadTimeout  time.Duration `default:"5s"       envconfig:"WRITE_TIMEOUT"`
	WriteTimeout time.Duration `default:"30s"      envconfig:"IDLE_TIMEOUT"`
	IdleTimeout  time.Duration `default:"30s"      envconfig:"IDLE_TIMEOUT"`
}

func main() {
	var config config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}

	storage := signature.NewMemory()

	router := http.NewServeMux()

	// Chain api documentation.
	router.Handle("/", docs.NewHandler())

	// Chain other services below...
	router.Handle("/signature/", http.StripPrefix("/signature", signature.NewHandler(storage)))
	// router.Handle("/account/", http.StripPrefix("/account", account.NewHandler())
	// router.Handle("/cart/", http.StripPrefix("/cart",  cart.NewHandler())
	// router.Handle("/wallet/", http.StripPrefix("/wallet",  wallet.NewHandler())
	// etc...

	// Chain k8s health checks, for now we do not have any external services but in future we can use errgroup and if any service above reports problem.
	router.HandleFunc("GET /liveness", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	router.HandleFunc("GET /readiness", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

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
