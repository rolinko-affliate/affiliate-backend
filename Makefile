.PHONY: build run test clean migrate-up migrate-down

# Build the application
build:
	go build -o affiliate-backend ./cmd/api

# Run the application
run:
	go run ./cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f affiliate-backend

# Run database migrations up
migrate-up:
	@echo "Running database migrations..."
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "DATABASE_URL is not set. Please set it and try again."; \
		exit 1; \
	fi
	migrate -path migrations -database "$(DATABASE_URL)" up

# Run database migrations down
migrate-down:
	@echo "Rolling back database migrations..."
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "DATABASE_URL is not set. Please set it and try again."; \
		exit 1; \
	fi
	migrate -path migrations -database "$(DATABASE_URL)" down

# Generate a new migration file
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "NAME is not set. Please set it and try again. Example: make migrate-create NAME=add_users_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations -seq $(NAME)

# Install dependencies
deps:
	go mod download

# Run linter
lint:
	go vet ./...
	golint ./...

# Generate a random encryption key
gen-key:
	@echo "Generating a random 32-byte base64 encoded encryption key..."
	@openssl rand -base64 32