package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/config"
	"github.com/affiliate-backend/internal/platform/logger"
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
		logger.Fatal("Unable to parse database URL", "error", err)
	}

	// Configure connection pool settings
	dbConfig.MaxConns = 25
	dbConfig.MinConns = 5
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = 30 * time.Minute

	logger.Debug("Database connection pool configuration", 
		"max_conns", dbConfig.MaxConns,
		"min_conns", dbConfig.MinConns,
		"max_conn_lifetime", dbConfig.MaxConnLifetime,
		"max_conn_idle_time", dbConfig.MaxConnIdleTime)

	// Create the connection pool
	DB, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		logger.Fatal("Unable to connect to database", "error", err)
	}

	// Ping the database to ensure connectivity
	err = DB.Ping(context.Background())
	if err != nil {
		logger.Fatal("Failed to ping database", "error", err)
	}

	logger.Info("Successfully connected to the database")
}

// CloseDB closes the database connection pool
func CloseDB() {
	if DB != nil {
		DB.Close()
		logger.Info("Database connection closed")
	}
}
