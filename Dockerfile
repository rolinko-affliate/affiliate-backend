# Build stage
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

# Build arguments for cross-compilation
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application and migration tool
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o affiliate-backend ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o migrate ./cmd/migrate

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binaries from the builder stage
COPY --from=builder /app/affiliate-backend .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations

# Install CA certificates for HTTPS connections
RUN apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates

# Set environment variables
ENV TZ=UTC \
    APP_ENV=production \
    ENVIRONMENT=production

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./affiliate-backend"]