# Signature service

- Explanation
- Run web API
- Running tests
- Running linter

## Explanation

- **`README.md`**: Contains project documentation, including this structure overview.
- **`cmd/`**: Holds executable commands:
  - `migrator/`: Command for handling database migrations.
  - `web/`: Command for starting the web server.
- **`go.mod`** & **`go.sum`**: Dependency management files.
- **`pkg/`**: Core business logic and package modules:
  - `cryptic/`: Handles cryptographic operations such as RSA and ECDSA.
  - `docs/`: Serves API documentation and related templates.
  - `migrator/`: Manages database migrations with SQL scripts.
  - `signature/`: Manages signature devices and transactions, including an in-memory storage implementation.
- **`task.md`**: Contains project-related tasks or requirements.

## Run web API

```sh
go mod tidy
go run cmd/web/main.go
```

## Running tests

```sh
go test ./...
```

## Running linter

```sh
golangci-lint run cmd/web
golangci-lint run cmd/migrator
golangci-lint run pkg/signature 
# etc ...
```
