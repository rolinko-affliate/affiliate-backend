.PHONY: build build-all run run-with-migrate test clean lint deps gen-key version \
	migrate-up migrate-down migrate-reset migrate-check migrate-version migrate-status migrate-create \
	docker-build docker-push docker-run docker-stop docker-migrate-up docker-migrate-down docker-migrate-check

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=affiliate-backend
API_BINARY=api
MIGRATE_BINARY=migrate
VERSION=$(shell cat VERSION 2>/dev/null || echo "dev")
IMAGE_NAME=asia-east2-docker.pkg.dev/jinko-test/jinko-test-docker-repo/saas-app

# Build the API application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/api

# Build all binaries
build-all:
	$(GOBUILD) -o $(API_BINARY) -v ./cmd/api
	$(GOBUILD) -o $(MIGRATE_BINARY) -v ./cmd/migrate

# Run the application
run:
	$(GORUN) ./cmd/api/main.go --mock-mode

# Run the application with auto-migrate
run-with-migrate:
	$(GORUN) ./cmd/api/main.go --auto-migrate

# Generate OpenAPI specification
swagger:
	$(shell go env GOPATH)/bin/swag init -g cmd/api/main.go -o docs

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(API_BINARY) $(MIGRATE_BINARY)

# Run linter
lint:
	go vet ./...
	$(shell go env GOPATH)/bin/golint ./...

# Install dependencies
deps:
	go mod download
	@echo "Installing required packages for migrations..."
	go get -u github.com/golang-migrate/migrate/v4
	go get -u github.com/golang-migrate/migrate/v4/database/postgres
	go get -u github.com/golang-migrate/migrate/v4/source/file
	@echo "Installing development tools..."
	go get -u github.com/jackc/pgx/v5
	go get -u github.com/swaggo/swag/cmd/swag
	go get -u golang.org/x/lint/golint

# Show version
version:
	@echo "Version: $(VERSION)"

# Generate a random encryption key
gen-key:
	@echo "Generating a random 32-byte base64 encoded encryption key..."
	@openssl rand -base64 32

# Database migration commands
# Run database migrations up
migrate-up:
	@echo "Running database migrations..."
	$(GORUN) ./cmd/migrate/main.go up

# Run database migrations down (rollback one migration)
migrate-down:
	@echo "Rolling back one migration..."
	$(GORUN) ./cmd/migrate/main.go down

# Run database migrations down (rollback all migrations)
migrate-reset:
	@echo "Rolling back all migrations..."
	$(GORUN) ./cmd/migrate/main.go reset

# Check database migration status
migrate-check:
	@echo "Checking migration status..."
	$(GORUN) ./cmd/migrate/main.go check

# Show current database version
migrate-version:
	@echo "Current database version:"
	$(GORUN) ./cmd/migrate/main.go version

# Show detailed migration status
migrate-status:
	@echo "Migration status:"
	$(GORUN) ./cmd/migrate/main.go status

# Generate a new migration file
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "NAME is not set. Please set it and try again. Example: make migrate-create NAME=add_users_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations -seq $(NAME)

# Docker commands
# Build Docker image for linux/amd64 (GKE compatible)
docker-build:
	docker build --platform linux/amd64 -t $(IMAGE_NAME):$(VERSION) .

# Build multi-platform Docker image (recommended for production)
docker-build-multi:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(IMAGE_NAME):$(VERSION) --push .

# Push Docker image
docker-push:
	docker push $(IMAGE_NAME):$(VERSION)

# Run Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker Compose
docker-stop:
	docker-compose down

# Run database migrations up using Docker
docker-migrate-up:
	docker-compose run --rm migrate up

# Run database migrations down using Docker
docker-migrate-down:
	docker-compose run --rm migrate down

# Check database migration status using Docker
docker-migrate-check:
	docker-compose run --rm migrate check

# Help
help:
	@echo "Available commands:"
	@echo "  make build              - Build the API application"
	@echo "  make build-all          - Build all binaries (API and migrate)"
	@echo "  make run                - Run the application"
	@echo "  make run-with-migrate   - Run the application with auto-migrate"
	@echo "  make test               - Run tests"
	@echo "  make clean              - Clean build files"
	@echo "  make lint               - Run linter"
	@echo "  make deps               - Install dependencies"
	@echo "  make version            - Show version"
	@echo "  make swagger            - Generate Swagger documentation"
	@echo "  make gen-key            - Generate a random encryption key"
	@echo "  make migrate-up         - Run database migrations up"
	@echo "  make migrate-down       - Rollback one migration"
	@echo "  make migrate-reset      - Rollback all migrations"
	@echo "  make migrate-check      - Check if migrations are up to date"
	@echo "  make migrate-version    - Show current database version"
	@echo "  make migrate-status     - Show detailed migration status"
	@echo "  make migrate-create     - Generate a new migration file"
	@echo "  make docker-build       - Build Docker image for linux/amd64"
	@echo "  make docker-build-multi - Build multi-platform Docker image and push"
	@echo "  make docker-push        - Push Docker image"
	@echo "  make docker-run         - Run Docker Compose"
	@echo "  make docker-stop        - Stop Docker Compose"
	@echo "  make docker-migrate-up  - Run migrations up using Docker"
	@echo "  make docker-migrate-down- Run migrations down using Docker"
	@echo "  make docker-migrate-check- Check migrations using Docker"