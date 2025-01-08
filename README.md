# Fuzzy Adventure

Go API platform with JWT auth, users/posts CRUD, PostgreSQL (GORM), Redis, Docker Compose, and unit tests.

## Features

- Register / login with JWT
- Posts CRUD for authenticated users
- Layered layout: handlers → services → repo
- Docker multi-stage build + compose
- `make test` coverage target

## Tech

Go · Gin · GORM · PostgreSQL · Redis · Docker

## Run

```bash
git clone https://github.com/d28035203/fuzzy-adventure.git
cd fuzzy-adventure
cp .env.example .env
docker compose up --build
# or: make run
```

## License

MIT
