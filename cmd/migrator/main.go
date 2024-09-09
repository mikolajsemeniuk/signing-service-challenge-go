package main

import (
	"fmt"
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/migrator"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

// This is an example relational database migrator for PostgreSQL that utilizes the `golang-migrate`
// library to handle schema migrations. The migrations are embedded within the binary using Go's
// embed package, and the configuration is read from environment variables. It can be used as an
// initContainer in a Kubernetes setup to automatically run database migrations when starting
// services that depend on the database. The migration files are located in the `migrations` folder
// and are applied to the PostgreSQL database defined by the environment variables.

type config struct {
	User     string `envconfig:"USER"     required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Host     string `envconfig:"HOST"     required:"true"`
	Database string `envconfig:"DATABASE" required:"true"`
	SSL      string `envconfig:"SSL"      required:"true"`
	Version  uint   `envconfig:"VERSION"  required:"true"`
}

func main() {
	var config config
	if err := envconfig.Process("PSQL", &config); err != nil {
		log.Fatal(err)
	}

	source := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", config.User, config.Password, config.Host, config.Database, config.SSL)
	if err := migrator.Migrate(source, config.Version); err != nil {
		log.Fatal(err)
	}
}
