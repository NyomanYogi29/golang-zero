# go-practice

My personal playground for learning Go. No framework, just the standard library, PostgreSQL, and Redis, structured the way I'd want a real backend to look. Well, soon I'll be using a framework but not now.

## Tech stack

- **Go** (`net/http`) The standard http library
- **PostgreSQL** via [`pgx`](https://github.com/jackc/pgx) Manual migrations
- **Redis** for caching and soon I wanna implement a queue and worker
- **Air** for live reload during development

## Project structure

```
cmd/            entry point (main.go equivalent, boots the server)
internal/
  config/       Postgres & Redis connection setup
  database/     migrations and schema definitions
  handler/      HTTP handlers
  service/      business logic
  repository/   data access layer
  model/        request/response and domain structs
  global/       shared helpers (JSON responses, utilities)
```

## Getting started

1. Install the dev-only CLI tools ([air](https://github.com/air-verse/air) for live reload, [swag](https://github.com/swaggo/swag) for generating the OpenAPI docs):

   ```bash
   make tools
   ```

2. Copy the environment file and fill in your own values:

   ```bash
   cp .env.example .env
   ```

3. Spin up Postgres and Redis:

   ```bash
   make docker-up
   ```

4. Run the server with live reload:

   ```bash
   make dev
   ```

The API will be available at `http://localhost:8080`. Database migrations run automatically on startup. Swagger UI is served at `http://localhost:8080/swagger/index.html`.

## Available commands

```bash
make dev          # run with live reload (air) — also regenerates Swagger docs on every reload
make swagger      # regenerate the Swagger/OpenAPI docs in ./docs from source annotations
make build        # regenerate Swagger docs, then build the binary to ./bin/server
make test         # run tests
make docker-up    # start Postgres & Redis containers
make docker-down  # stop them
make tools        # install air and swag CLI tools
```

## Endpoints

| Method | Path                | Description       |
| ------ | ------------------- | ------------------ |
| GET    | `/`                 | Health check        |
| POST   | `/api/v1/register`  | Register a new user |
| POST   | `/api/v1/login`     | Log in and receive a JWT |
| POST   | `/api/v1/logout`    | Log out (requires `Authorization: Bearer <token>`) |

More will show up here as I keep experimenting.

## Adding a new endpoint

1. Add the handler in `internal/handler/`, with `swag` annotations above the function (`@Summary`, `@Param`, `@Success`, `@Failure`, `@Router`, ...) following the existing handlers as examples.
2. Register the route in `cmd/server.go`.
3. That's it — `make dev` regenerates `docs/` on every live-reload rebuild, so the new route shows up in Swagger UI as soon as you save. Run `make swagger` directly if you want to regenerate without starting the server.

Thankyou!
