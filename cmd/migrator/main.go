package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/migrator"
)

// This is an example relational database migrator for PostgreSQL that utilizes the `golang-migrate`
// library to handle schema migrations. The migrations are embedded within the binary using Go's
// embed package, and the configuration is read from environment variables. It can be used as an
// initContainer in a Kubernetes setup to automatically run database migrations when starting
// services that depend on the database. The migration files are located in the `migrations` folder
// and are applied to the PostgreSQL database defined by the environment variables.

func main() {
	if err := migrator.Migrate(); err != nil {
		log.Fatal(err)
	}
}
