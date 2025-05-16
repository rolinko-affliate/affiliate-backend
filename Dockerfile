# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o affiliate-backend ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/affiliate-backend .
COPY --from=builder /app/migrations ./migrations

# Install migrate for database migrations
RUN apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    apk del curl

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./affiliate-backend"]