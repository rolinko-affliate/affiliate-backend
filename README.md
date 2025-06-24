# Affiliate Backend Platform

A comprehensive affiliate marketing platform backend built with Go, featuring clean architecture, multi-tenant support, and external provider integrations.

## 🏗️ Architecture Overview

This project follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────────┐
│                        API Layer (Gin)                         │
├─────────────────────────────────────────────────────────────────┤
│  Handlers │ Middleware │ Router │ Models │ Error Handling      │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Service Layer                             │
├─────────────────────────────────────────────────────────────────┤
│  Business Logic │ Validation │ Orchestration │ Provider Sync   │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Domain Layer                               │
├─────────────────────────────────────────────────────────────────┤
│  Entities │ Value Objects │ Domain Events │ Business Rules     │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Repository Layer                             │
├─────────────────────────────────────────────────────────────────┤
│  Data Access │ Persistence │ Query Building │ Transactions     │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Provider Layer                               │
├─────────────────────────────────────────────────────────────────┤
│  Integration Service │ Everflow Client │ Mappers │ Mock Service │
└─────────────────────────────────────────────────────────────────┘
```

## 🚀 Key Features

### Core Functionality
- **Multi-tenant Organization Management**: Support for multiple organizations with isolated data
- **Advertiser Management**: Complete advertiser lifecycle with external provider sync
- **Affiliate Management**: Comprehensive affiliate partner management
- **Campaign Management**: Campaign creation, management, and tracking
- **Tracking Link Generation**: Dynamic tracking link creation with QR codes
- **Analytics & Reporting**: Performance metrics and data insights

### Technical Features
- **Clean Architecture**: Modular, testable, and maintainable codebase
- **Provider Integration**: Pluggable external provider system (Everflow)
- **Role-Based Access Control (RBAC)**: Fine-grained permission system
- **JWT Authentication**: Secure authentication with Supabase integration
- **Database Migrations**: Version-controlled schema management
- **Mock Mode**: Development-friendly mock integrations
- **API Documentation**: Auto-generated Swagger/OpenAPI docs

## 🛠️ Technology Stack

- **Language**: Go 1.23
- **Web Framework**: Gin
- **Database**: PostgreSQL with pgx driver
- **Authentication**: JWT with Supabase
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose
- **Build System**: Make
- **Testing**: Go testing framework with mocks

## 📁 Project Structure

```
├── cmd/
│   ├── api/           # Main API application
│   └── migrate/       # Database migration tool
├── internal/
│   ├── api/           # HTTP handlers, middleware, routing
│   ├── auth/          # Authentication logic
│   ├── config/        # Configuration management
│   ├── domain/        # Domain entities and business rules
│   ├── platform/      # External integrations (Everflow, crypto)
│   ├── repository/    # Data access layer
│   └── service/       # Business logic layer
├── migrations/        # Database migration files
├── docs/             # Auto-generated API documentation
├── k8s/              # Kubernetes deployment manifests
└── scripts/          # Utility scripts
```

## 🔧 Local Development Setup

### Prerequisites

- **Go 1.23+**
- **Docker & Docker Compose**
- **PostgreSQL** (or use Docker)
- **Make** (for build commands)

### Environment Variables

Create a `.env` file in the project root:

```bash
# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5432/affiliate_platform?sslmode=disable

# Or use individual components:
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=affiliate_platform
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_SSL_MODE=disable

# Authentication
SUPABASE_JWT_SECRET=your-supabase-jwt-secret-here

# Encryption (generate with: make gen-key)
ENCRYPTION_KEY=your-32-byte-base64-encoded-key-here

# Application Settings
PORT=8080
ENVIRONMENT=development
DEBUG_MODE=true
MOCK_MODE=true  # Use mock integrations for development
```

### Quick Start with Docker

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd affiliate-backend
   ```

2. **Start with Docker Compose**:
   ```bash
   # Start database and run migrations
   docker-compose up -d db
   docker-compose run --rm migrate-up
   
   # Start the application
   docker-compose up app
   ```

3. **Access the application**:
   - API: http://localhost:8080
   - Health Check: http://localhost:8080/health
   - API Documentation: http://localhost:8080/swagger/index.html

### Manual Setup

1. **Install dependencies**:
   ```bash
   make deps
   ```

2. **Generate encryption key**:
   ```bash
   make gen-key
   # Copy the output to your .env file as ENCRYPTION_KEY
   ```

3. **Start PostgreSQL** (if not using Docker):
   ```bash
   # Using Docker for just the database
   docker run -d \
     --name postgres-affiliate \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=affiliate_platform \
     -p 5432:5432 \
     postgres:14-alpine
   ```

4. **Run database migrations**:
   ```bash
   make migrate-up
   ```

5. **Start the application**:
   ```bash
   # Development mode with mock integrations
   make run
   
   # Or with auto-migration
   make run-with-migrate
   ```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test -v ./internal/service/...

# Run tests with coverage
go test -cover ./...
```

### Test Categories

1. **Unit Tests**: Test individual components in isolation
2. **Integration Tests**: Test component interactions
3. **Repository Tests**: Test database operations
4. **Service Tests**: Test business logic

### Mock Mode

The application supports mock mode for development and testing:

```bash
# Start in mock mode (default for make run)
go run ./cmd/api/main.go --mock-mode

# Or set environment variable
MOCK_MODE=true go run ./cmd/api/main.go
```

Mock mode replaces external provider integrations with logging mock services.

### Local Setup Validation

Use the included test script to validate your local development setup:

```bash
# Run the setup validation script
./test_local_setup.sh
```

This script will:
- ✅ Check required tools (Go, Make, Docker)
- ✅ Validate project structure
- ✅ Test Makefile commands
- ✅ Verify documentation exists
- ✅ Test environment configuration
- ✅ Validate application startup

## 📚 API Documentation

### Swagger/OpenAPI

1. **Generate documentation**:
   ```bash
   make swagger
   ```

2. **Access documentation**:
   - Start the server: `make run`
   - Visit: http://localhost:8080/swagger/index.html

### Key API Endpoints

#### Authentication
- `POST /api/v1/public/webhooks/supabase/new-user` - Supabase user webhook
- `GET /api/v1/users/me` - Get current user profile

#### Organizations
- `GET /api/v1/organizations` - List organizations
- `POST /api/v1/organizations` - Create organization (Admin only)
- `GET /api/v1/organizations/:id` - Get organization details

#### Advertisers
- `POST /api/v1/advertisers` - Create advertiser
- `GET /api/v1/advertisers/:id` - Get advertiser
- `PUT /api/v1/advertisers/:id` - Update advertiser
- `POST /api/v1/advertisers/:id/sync-to-everflow` - Sync to Everflow

#### Affiliates
- `POST /api/v1/affiliates` - Create affiliate
- `GET /api/v1/affiliates/:id` - Get affiliate
- `PUT /api/v1/affiliates/:id` - Update affiliate

#### Campaigns
- `POST /api/v1/campaigns` - Create campaign
- `GET /api/v1/campaigns/:id` - Get campaign
- `PUT /api/v1/campaigns/:id` - Update campaign

#### Tracking Links
- `POST /api/v1/organizations/:id/tracking-links` - Create tracking link
- `GET /api/v1/organizations/:id/tracking-links` - List tracking links
- `POST /api/v1/organizations/:id/tracking-links/generate` - Generate tracking link

#### Analytics
- `GET /api/v1/analytics/advertisers/:id` - Get advertiser analytics
- `GET /api/v1/analytics/affiliates/:id` - Get affiliate analytics

## 🔐 Authentication & Authorization

### JWT Authentication

The platform uses JWT tokens from Supabase for authentication:

1. **Token Validation**: All protected routes require a valid JWT token
2. **User Context**: User information is extracted from the token
3. **Organization Access**: Users are associated with organizations

### Role-Based Access Control (RBAC)

Supported roles:
- **Admin**: Full system access
- **AdvertiserManager**: Manage advertisers and campaigns
- **AffiliateManager**: Manage affiliates
- **User**: Read-only access to assigned resources

### Example Authentication

```bash
# Include JWT token in requests
curl -H "Authorization: Bearer your-jwt-token" \
     http://localhost:8080/api/v1/users/me
```

## 🗄️ Database Management

### Migrations

```bash
# Run migrations up
make migrate-up

# Rollback one migration
make migrate-down

# Check migration status
make migrate-status

# Create new migration
make migrate-create NAME=add_new_feature

# Reset all migrations (development only)
make migrate-reset
```

### Database Schema

Key tables:
- `organizations` - Multi-tenant organization data
- `profiles` - User profiles linked to Supabase
- `advertisers` - Advertiser entities
- `affiliates` - Affiliate partner entities
- `campaigns` - Marketing campaigns
- `tracking_links` - Generated tracking links
- `*_provider_mappings` - External provider relationships

## 🔌 Provider Integration

### Everflow Integration

The platform integrates with Everflow for affiliate tracking:

1. **Provider Mapping Pattern**: Separate entities for provider-specific data
2. **Sync Operations**: Bi-directional synchronization with Everflow
3. **Data Transformation**: Clean mapping between internal and external models

### Adding New Providers

1. Implement the `IntegrationService` interface
2. Create provider-specific mappers
3. Add provider mapping repositories
4. Update configuration

## 🚀 Deployment

### Docker Deployment

```bash
# Build Docker image
make docker-build

# Run with Docker Compose
make docker-run

# Stop services
make docker-stop
```

### Kubernetes Deployment

Kubernetes manifests are available in the `k8s/` directory:

```bash
# Apply base configuration
kubectl apply -k k8s/base/

# Apply environment-specific overlays
kubectl apply -k k8s/overlays/development/
```

### Environment Configuration

For production deployment, ensure these environment variables are set:

- `DATABASE_URL` - PostgreSQL connection string
- `SUPABASE_JWT_SECRET` - Supabase JWT secret
- `ENCRYPTION_KEY` - 32-byte encryption key for sensitive data
- `ENVIRONMENT=production`
- `MOCK_MODE=false`

## 🔧 Development Commands

```bash
# Build application
make build

# Build all binaries
make build-all

# Run application (mock mode)
make run

# Run with auto-migration
make run-with-migrate

# Run tests
make test

# Generate API documentation
make swagger

# Clean build artifacts
make clean

# Install dependencies
make deps

# Generate encryption key
make gen-key

# Show version
make version

# Show all available commands
make help
```

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**:
   - Ensure PostgreSQL is running
   - Check `DATABASE_URL` configuration
   - Verify database exists and is accessible

2. **Migration Errors**:
   - Check migration status: `make migrate-status`
   - Ensure database is clean for initial setup
   - Use `make migrate-reset` for development (destructive)

3. **Authentication Errors**:
   - Verify `SUPABASE_JWT_SECRET` is correct
   - Check JWT token format and expiration
   - Ensure user has proper role assignments

4. **Provider Integration Issues**:
   - Use mock mode for development: `--mock-mode`
   - Check provider credentials and configuration
   - Review provider mapping data

### Debug Mode

Enable debug mode for detailed logging:

```bash
DEBUG_MODE=true go run ./cmd/api/main.go
```

### Health Checks

Monitor application health:

```bash
# Basic health check
curl http://localhost:8080/health

# Database connectivity check
curl http://localhost:8080/api/v1/health/db
```

## 📖 Additional Documentation

- [Domain Design Pattern](./domain_design_pattern.md) - Detailed architecture explanation
- [API Documentation](./docs/) - Auto-generated Swagger docs
- [Migration Guide](./migrations/MIGRATION_UPDATE_SUMMARY.md) - Database migration details
- [Provider Integration](./internal/platform/provider/README.md) - External provider setup

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `make test`
6. Submit a pull request

## 🚀 Quick Reference

### Essential Commands

```bash
# Setup
make gen-key                    # Generate encryption key
make deps                       # Install dependencies

# Development
make run                        # Start with mock mode
make run-with-migrate          # Start with auto-migrate
make test                      # Run all tests
make swagger                   # Generate API docs

# Database
make migrate-up                # Apply migrations
make migrate-down              # Rollback one migration
make migrate-status            # Check migration status

# Docker
docker-compose up -d db        # Start database only
docker-compose up              # Start all services
docker-compose run --rm migrate-up  # Run migrations in Docker

# Validation
./test_local_setup.sh          # Validate local setup
```

### Quick Start Checklist

1. ✅ Clone repository
2. ✅ Install Go 1.23+ and Make
3. ✅ Create `.env` file with required variables
4. ✅ Generate encryption key: `make gen-key`
5. ✅ Start database: `docker-compose up -d db`
6. ✅ Run migrations: `make migrate-up`
7. ✅ Start application: `make run`
8. ✅ Test setup: `./test_local_setup.sh`
9. ✅ Access API docs: http://localhost:8080/swagger/index.html

## 📄 License

This project is licensed under the Apache 2.0 License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Check the troubleshooting section above
- Review the additional documentation
- Create an issue in the repository
- Contact the development team

---

**Version**: 0.0.6
**Last Updated**: 2025-06-24