package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/affiliate-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the global database connection pool
var DB *pgxpool.Pool

// InitDBConnection creates a new database connection pool for a specific purpose
// This is useful for one-off operations like checking migration status
func InitDBConnection(dbURL string) (*pgxpool.Pool, error) {
	// Parse the database connection string
	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}
	
	// Configure connection pool settings for temporary use
	dbConfig.MaxConns = 5
	dbConfig.MinConns = 1
	dbConfig.MaxConnLifetime = 5 * time.Minute
	dbConfig.MaxConnIdleTime = 1 * time.Minute

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	// Ping the database to ensure connectivity
	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return pool, nil
}

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