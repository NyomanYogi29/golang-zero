# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A personal learning project for Go backend development. Standard library `net/http` only (no web framework), PostgreSQL via `pgx`, Redis for session storage, structured like a real production backend.

## Commands

```bash
make dev          # run with live reload (air) — primary way to run during development
make build        # build binary to ./bin/server
make test         # go test -v ./internal/...
make docker-up    # start Postgres & Redis containers (docker-compose)
make docker-down  # stop containers
```

Run a single test: `go test -v ./internal/<package> -run TestName`

Setup: copy `.env.example` to `.env` and fill in values, then `make docker-up` before `make dev`. Migrations run automatically on server startup — there is no separate migrate command.

No linter is configured in this repo.

## Architecture

Layered structure, one feature ("user") currently implemented end to end — follow this same layering for any new feature:

```
handler → service → repository → model
```

- `cmd/server.go` — entry point. Loads `.env` (via godotenv), connects Postgres/Redis, runs migrations, wires up repository → service → handler → middleware, and registers routes on a plain `http.ServeMux`. There is no router framework; new routes are added directly here with `mux.HandleFunc`.
- `internal/config/` — Postgres (`pgxpool`) and Redis client singletons, accessed via `GetDBPool()` / `GetRedisClient()` package-level getters rather than dependency injection. Note: `ConnectPostgres` currently reads `DATABASE_UTL` (typo) instead of `DATABASE_URL` from `.env.example`, so it silently falls back to the hardcoded local connection string — check this if Postgres connection env vars appear to have no effect.
- `internal/database/` — raw SQL schema strings (e.g. `user_schema.go`) collected in `allSchemas` and applied transactionally by `RunMigrations` on every startup. `CREATE TABLE IF NOT EXISTS` makes this idempotent. To add a table, write a new schema constant and append it to `allSchemas` in `migration.go`.
- `internal/repository/` — defines a `UserRepository` interface (`user_repository.go`) plus a concrete `PostgresUserRepository` (`user_postgres_repository.go`) using raw SQL via `pgx`, no ORM. Services depend on the interface, not the concrete type, so they can be tested with a mock/fake.
- `internal/service/` — business logic: password hashing (bcrypt), JWT issuance, and session bookkeeping. Sessions are stored in Redis as `session:<userID>` → token string, with TTL matching JWT expiry (24h). Login/logout/session-validity all go through Redis, not just JWT verification — a valid JWT alone does not authorize a request if the Redis session was deleted (logout) or overwritten (re-login elsewhere invalidates the old token).
- `internal/handler/` — plain `http.HandlerFunc` methods; manually checks `r.Method` since `ServeMux` here doesn't do method-based routing. Decodes JSON body into `model.*RequestSchema` structs, delegates to service, wraps results with `internal/global` response helpers.
- `internal/middleware/` — `AuthMiddleware.Authenticate` wraps a handler: parses `Bearer <token>`, validates the JWT, then re-checks session validity against Redis, then injects `user_id` into the request context (via bare string key `"user_id"`, not a typed context key).
- `internal/model/` — domain struct (`User`) plus separate request/response DTOs per endpoint (e.g. `UserRegisterRequestSchema`, `UserLoginResponseSchema`). Handlers/services should keep using these DTOs rather than exposing `User` directly (note `Register` zeroes `user.Password` before responding).
- `internal/global/` — `WriteJSON` plus `NewSuccessResponse`/`NewErrorResponse` builders used by every handler for a consistent `{success, message, data|error}` envelope.
- `internal/utils/jwt.go` — HS256 JWT generation/validation; secret from `JWT_SECRET_KEY` env var (has an insecure hardcoded fallback if unset — always set this in `.env`).

## API docs

Swagger annotations live above `RunServer()` in `cmd/server.go` and are used by `swaggo` to generate `docs/docs.go` / `docs/swagger.json` / `docs/swagger.yaml`, served at `/swagger/*`. Regenerate with `swag init` after changing route annotations (swag is not vendored as a Makefile target here, so run it directly if installed).

## Conventions to follow

- New endpoints: add route in `cmd/server.go`, request/response structs in `internal/model`, business logic in `internal/service`, DB access behind a repository interface in `internal/repository`.
- Error responses always go through `global.NewErrorResponse` with an UPPER_SNAKE_CASE error code; success responses through `global.NewSuccessResponse`.
- Some user-facing error/log strings are in Indonesian, some in English — the codebase is mixed; match whichever the surrounding file/function already uses rather than converting.
