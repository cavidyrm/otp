.PHONY: help build run test clean docker-build docker-run docker-stop docker-logs swagger

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker services"
	@echo "  docker-logs  - View Docker logs"
	@echo "  swagger      - Generate Swagger documentation"

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application locally
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Build Docker image
docker-build:
	docker build -t otp-service .

# Run with Docker Compose
docker-run:
	docker compose up -d

# Stop Docker services
docker-stop:
	docker compose down

# View Docker logs
docker-logs:
	docker compose logs -f

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Database migrations (local)
migrate-local:
	@echo "Make sure PostgreSQL is running and configured in .env"
	go run cmd/server/main.go

# Health check
health:
	curl -f http://localhost:8080/health || echo "Service is not running"

# API documentation
docs:
	@echo "Swagger documentation available at: http://localhost:8080/swagger/index.html"
