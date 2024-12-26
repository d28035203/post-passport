.PHONY: run build test clean lint fmt deps deps-test test test-coverage test-integration docker-up docker-down docker-logs test-unit

# Build
build:
	go build -o bin/api ./cmd/api

# Run
run: build
	./bin/api

# Development (hot-reload with air)
dev: deps
	go run -mod=mod ./cmd/api

# Test - Run all unit tests
test: deps-test
	go test -v -cover -race ./tests/

# Coverage - Run all tests with coverage
test-coverage: deps-test
	go test -v -coverprofile=coverage.out -covermode=atomic ./tests/
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html
	@echo "Coverage report generated: coverage.html"

# Test with Docker (integration tests)
test-integration:
	@echo "Running integration tests..."
	docker-compose exec app go test -v -cover -race github.com/d28035203/fuzzy-adventure/tests/ | tee test-integration-output.log

# Dependencies
deps:
	go mod download

# Test dependencies
deps-test:
	go get -u github.com/stretchr/testify@v1.9.0
	go mod tidy

# Clean
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache -modcache

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Quick test
test-quick: deps-test
	go test -v ./tests/

# Full test suite
test-full: deps-test
	@echo "Running full test suite..."
	go test -v -cover -race ./tests/
	@echo "✅ All tests passed!"
