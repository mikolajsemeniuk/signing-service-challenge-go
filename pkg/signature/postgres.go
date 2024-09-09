package signature

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
This is an example implementation of a PostgreSQL integration for the `signature` package.
The `Postgres` struct provides basic methods for connecting to a PostgreSQL database,
retrieving devices (via `ListDevices` and `FindDevice`), and managing database connections.
It showcases how you can implement persistent storage for the signature API using SQL queries with `database/sql` in Go.
*/

// Ensures interface is implement for proof of concept.
var _ Storage = &Postgres{}

// NewPostgres initializes a new connection to the PostgreSQL database.
func NewPostgres(p *pgxpool.Pool) *Postgres {
	return &Postgres{pool: p}
}

// Postgres represents a basic integration with a PostgreSQL database.
type Postgres struct {
	pool *pgxpool.Pool
}

// ListDevices implements Storage.
func (p *Postgres) ListDevices(_ context.Context) ([]Device, error) {
	panic("unimplemented")
}

// FindDevice implements Storage.
func (p *Postgres) FindDevice(_ context.Context, _ uuid.UUID) (Device, error) {
	panic("unimplemented")
}

// CreateDevice implements Storage.
func (p *Postgres) CreateDevice(_ context.Context, _ CreateDeviceInput) (Device, error) {
	panic("unimplemented")
}

// CreateTransaction implements Storage.
func (p *Postgres) CreateTransaction(_ context.Context, _ CreateTransactionInput) (Transaction, error) {
	panic("unimplemented")
}
