# Go Users & Posts API

> Repository: `post-passport` · https://github.com/d28035203/post-passport

Production-shaped **Go REST API** with JWT authentication, **users/posts CRUD**, layered handlers → services → repository design, **PostgreSQL** (GORM), **Redis** config hooks, multi-stage **Docker** + Compose, and unit tests with coverage.

## Features

- Register / login with JWT
- Authenticated profile endpoint
- Posts CRUD (list, get, create, update, delete)
- Layered layout under `internal/` (handlers → services → repo)
- Middleware: CORS, recover, request logging, metrics, JWT auth
- GORM auto-migrate for `User` and `Post`
- Multi-stage Docker image + Compose for local full stack
- `make test` / `make test-coverage` for the API test suite

## Tech stack

| Layer | Choice |
|-------|--------|
| Language | Go 1.21+ |
| HTTP | Gin |
| ORM / DB | GORM + PostgreSQL (SQLite available in module for tests) |
| Cache / config | Redis host settings in env (Compose) |
| Auth | JWT (`golang-jwt`) |
| Ops | Docker, Docker Compose, Makefile |

## Architecture

```
cmd/api                HTTP entrypoint (Gin router, wiring)
internal/
  handlers/            HTTP adapters
  services/            Business logic
  repo/                GORM data access
  middleware/          CORS, JWT, logging, metrics, recover
  models/              User, Post
  config/              Env + DB connection
pkg/logger             Shared logging helper
tests/                 API tests (testify)
```

Request flow: **Router → middleware → handler → service → repo → PostgreSQL**.

## API

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | No | Liveness (`status`, `service`) |
| POST | `/register` | No | Create user |
| POST | `/login` | No | Issue JWT |
| GET | `/profile` | JWT | Current user |
| GET | `/posts` | JWT | List posts |
| GET | `/posts/:id` | JWT | Get post |
| POST | `/posts` | JWT | Create post |
| PUT | `/posts/:id` | JWT | Update post |
| DELETE | `/posts/:id` | JWT | Delete post |

## Prerequisites

- Go 1.21+
- Docker + Docker Compose (recommended local path)
- Or local PostgreSQL (+ optional Redis) matching `.env.example`

## Run with Docker

```bash
git clone https://github.com/d28035203/post-passport.git
cd post-passport
cp .env.example .env
docker compose up --build
```

API defaults to **http://localhost:8080** (see `PORT` in `.env`).

Health check:

```bash
curl -s http://localhost:8080/health
```

## Run locally (Go)

```bash
cp .env.example .env
# start Postgres (and Redis if you use it) — or use: make docker-up
make deps
make run
# or: make dev
```

## Tests

```bash
make test              # unit/API tests with race + coverage flags
make test-coverage     # HTML coverage report
make test-quick        # faster pass
```

## Makefile targets

| Target | Description |
|--------|-------------|
| `make run` / `make dev` | Build/run the API |
| `make test` | Tests with race detector |
| `make test-coverage` | Coverage HTML |
| `make docker-up` / `docker-down` | Compose helpers |
| `make lint` / `make fmt` | Code hygiene (if configured in Makefile) |
| `make clean` | Remove `bin/` and coverage artifacts |

## Structure

```
post-passport/
├── cmd/api/main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repo/
│   └── services/
├── pkg/logger/
├── tests/
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .env.example
└── README.md
```

## Configuration

Copy `.env.example` → `.env`. Important keys:

| Variable | Purpose |
|----------|---------|
| `DB_*` | PostgreSQL connection |
| `REDIS_*` | Redis host/port (Compose stack) |
| `JWT_SECRET` / `JWT_EXP` | Token signing and expiry |
| `PORT` | HTTP listen port (default `8080`) |
| `LOG_LEVEL` | Logger verbosity |

Change `JWT_SECRET` before any real deployment.

## License

MIT
