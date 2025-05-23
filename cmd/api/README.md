# API Command

This module contains the main entry point for the API server application. It initializes all components, sets up the HTTP server, and handles graceful shutdown.

## Key Components

### Main Function

The `main` function is the entry point of the application:

```go
func main() {
    // Load Configuration
    config.LoadConfig()
    appConf := config.AppConfig

    // Check database migration status
    if err := checkDatabaseMigrations(&appConf, autoMigrate); err != nil {
        log.Fatalf("Database migration check failed: %v", err)
    }

    // Initialize Database
    repository.InitDB(&appConf)
    defer repository.CloseDB()

    // Initialize Repositories, Services, and Handlers
    // ...

    // Setup Router
    router := api.SetupRouter(api.RouterOptions{
        // ...
    })

    // Start Server
    srv := &http.Server{
        Addr:    ":" + appConf.Port,
        Handler: router,
    }

    // Start the server in a goroutine
    go func() {
        log.Printf("Server starting on port %s\n", appConf.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("ListenAndServe(): %s\n", err)
        }
    }()

    // Graceful Shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exiting")
}
```

### Database Migration Check

The `checkDatabaseMigrations` function verifies that the database schema is up to date:

```go
func checkDatabaseMigrations(cfg *config.Config, autoMigrate bool) error {
    // Implementation details...
}
```

## Initialization Process

The application follows a specific initialization sequence:

1. Load configuration from environment variables
2. Check database migration status
3. Initialize database connection
4. Initialize repositories
5. Initialize platform services (crypto, everflow)
6. Initialize domain services
7. Initialize API handlers
8. Setup router with handlers and middleware
9. Start HTTP server
10. Setup graceful shutdown

## Component Wiring

The application wires together all components:

```go
// Initialize Repositories
profileRepo := repository.NewPgxProfileRepository(repository.DB)
organizationRepo := repository.NewPgxOrganizationRepository(repository.DB)
advertiserRepo := repository.NewPgxAdvertiserRepository(repository.DB)
campaignRepo := repository.NewPgxCampaignRepository(repository.DB)
affiliateRepo := repository.NewPgxAffiliateRepository(repository.DB)

// Initialize Platform Services
cryptoService := crypto.NewServiceFromConfig()
everflowService, err := everflow.NewEverflowServiceFromEnv(advertiserRepo, campaignRepo, cryptoService)

// Initialize Domain Services
profileService := service.NewProfileService(profileRepo)
organizationService := service.NewOrganizationService(organizationRepo)
advertiserService := service.NewAdvertiserService(advertiserRepo, organizationRepo, everflowService, cryptoService)
affiliateService := service.NewAffiliateService(affiliateRepo, organizationRepo)
campaignService := service.NewCampaignService(campaignRepo, advertiserRepo, organizationRepo, everflowService, cryptoService)

// Initialize Handlers
profileHandler := handlers.NewProfileHandler(profileService)
organizationHandler := handlers.NewOrganizationHandler(organizationService, profileService)
advertiserHandler := handlers.NewAdvertiserHandler(advertiserService, profileService)
affiliateHandler := handlers.NewAffiliateHandler(affiliateService, profileService)
campaignHandler := handlers.NewCampaignHandler(campaignService)
```

## Command-Line Flags

The application supports command-line flags:

- `--auto-migrate`: Automatically apply pending database migrations

## Graceful Shutdown

The application implements graceful shutdown:

1. Capture termination signals (SIGINT, SIGTERM)
2. Stop accepting new requests
3. Wait for ongoing requests to complete (with timeout)
4. Close database connections
5. Exit the application

## Error Handling

The application handles errors at the top level:

- Fatal errors during initialization stop the application
- Non-fatal errors are logged and handled appropriately
- Panic recovery middleware catches unexpected panics