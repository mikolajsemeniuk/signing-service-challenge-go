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
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(source string, version uint) error {
	driver, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("error creating migration driver: %w", err)
	}

	migration, err := migrate.NewWithSourceInstance("iofs", driver, source)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %w", err)
	}
	defer migration.Close()

	if err = migration.Migrate(version); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migration: %w", err)
	}

	return nil
}
