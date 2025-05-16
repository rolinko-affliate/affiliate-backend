package repository

import (
	"context"
	"log"
	"time"

	"github.com/affiliate-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the global database connection pool
var DB *pgxpool.Pool

// InitDB initializes the database connection pool
func InitDB(cfg *config.Config) {
	var err error
	// Parse the database connection string
	dbConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}
	
	// Configure connection pool settings
	dbConfig.MaxConns = 25
	dbConfig.MinConns = 5
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = 30 * time.Minute

	// Create the connection pool
	DB, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Ping the database to ensure connectivity
	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
	}

	log.Println("Successfully connected to the database!")
}

// CloseDB closes the database connection pool
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed.")
	}
}