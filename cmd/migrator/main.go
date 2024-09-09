package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

const Version = 1

//go:embed migrations/*.sql
var migrations embed.FS

type config struct {
	User     string `envconfig:"PSQL_USER"     required:"true"`
	Password string `envconfig:"PSQL_PASSWORD" required:"true"`
	Host     string `envconfig:"PSQL_HOST"     required:"true"`
	Database string `envconfig:"PSQL_DATABASE" required:"true"`
	SSL      string `envconfig:"PSQL_SSL"      required:"true"`
}

func main() {
	var config config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}

	driver, err := iofs.New(migrations, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	source := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", config.User, config.Password, config.Host, config.Database, config.SSL)
	migration, err := migrate.NewWithSourceInstance("iofs", driver, source)
	if err != nil {
		log.Fatal(err)
	}

	defer migration.Close()
	if err = migration.Migrate(Version); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
