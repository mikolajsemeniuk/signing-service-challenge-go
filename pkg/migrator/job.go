package migrator

// This is an example relational database migrator for PostgreSQL that utilizes the `golang-migrate`
// library to handle schema migrations. The migrations are embedded within the binary using Go's
// embed package, and the configuration is read from environment variables. It can be used as an
// initContainer in a Kubernetes setup to automatically run database migrations when starting
// services that depend on the database. The migration files are located in the `migrations` folder
// and are applied to the PostgreSQL database defined by the environment variables.

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

//go:embed migrations/*.sql
var migrations embed.FS

type config struct {
	User     string `envconfig:"USER"     required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Host     string `envconfig:"HOST"     required:"true"`
	Database string `envconfig:"DATABASE" required:"true"`
	SSL      string `envconfig:"SSL"      required:"true"`
	Version  uint   `envconfig:"VERSION"  required:"true"`
}

func Migrate() error {
	var config config
	if err := envconfig.Process("PSQL", &config); err != nil {
		return err
	}

	driver, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	source := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", config.User, config.Password, config.Host, config.Database, config.SSL)

	migration, err := migrate.NewWithSourceInstance("iofs", driver, source)
	if err != nil {
		return err
	}
	defer migration.Close()

	if err = migration.Migrate(config.Version); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
