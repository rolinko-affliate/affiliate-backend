package main

// @title           Affiliate Backend API
// @version         1.0
// @description     API Server for Affiliate Backend Application
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/affiliate-backend/internal/api"
	"github.com/affiliate-backend/internal/api/handlers"
	"github.com/affiliate-backend/internal/config"
	"github.com/affiliate-backend/internal/repository"
	"github.com/affiliate-backend/internal/service"
)

func main() {
	// Load Configuration
	config.LoadConfig()
	appConf := config.AppConfig

	// Initialize Database
	repository.InitDB(&appConf)
	defer repository.CloseDB()

	// Initialize Repositories
	profileRepo := repository.NewPgxProfileRepository(repository.DB)
	// Initialize other repositories as needed

	// Initialize Services
	profileService := service.NewProfileService(profileRepo)
	// Initialize other services as needed

	// Initialize Handlers
	profileHandler := handlers.NewProfileHandler(profileService)
	// Initialize other handlers as needed

	// Setup Router
	router := api.SetupRouter(api.RouterOptions{
		ProfileHandler: profileHandler,
		ProfileService: profileService,
		// Add other handlers and services as needed
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