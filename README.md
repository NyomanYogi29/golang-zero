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

1. Copy the environment file and fill in your own values:

   ```bash
   cp .env.example .env
   ```

2. Spin up Postgres and Redis:

   ```bash
   make docker-up
   ```

3. Run the server with live reload:

   ```bash
   make dev
   ```

The API will be available at `http://localhost:8080`. Database migrations run automatically on startup.

## Available commands

```bash
make dev          # run with live reload (air)
make build        # build the binary to ./bin/server
make test         # run tests
make docker-up    # start Postgres & Redis containers
make docker-down  # stop them
```

## Endpoints

| Method | Path                | Description       |
| ------ | ------------------- | ------------------ |
| GET    | `/`                 | Health check        |
| POST   | `/api/v1/register`  | Register a new user |

More will show up here as I keep experimenting.
Thankyou!
