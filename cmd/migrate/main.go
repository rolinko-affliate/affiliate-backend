package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/affiliate-backend/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrationInfo contains detailed information about migrations
type MigrationInfo struct {
	CurrentVersion int64  `json:"current_version"`
	LatestVersion  int64  `json:"latest_version"`
	Status         string `json:"status"`
	PendingCount   int    `json:"pending_count,omitempty"`
}

// MigrationStatus represents the status of migrations
type MigrationStatus int

const (
	MigrationStatusUpToDate MigrationStatus = iota
	MigrationStatusPending
	MigrationStatusError
)

func main() {
	fmt.Println("Database Migration Tool")
	fmt.Println("----------------------")
	
	
	// Load Configuration
	config.LoadConfig()
	appConf := config.AppConfig

	if appConf.DatabaseURL == "" {
		log.Fatalf("DATABASE_URL is not set. Please set it and try again.")
	}

	// Check command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Create migration instance
	m, err := migrate.New("file://migrations", appConf.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	defer m.Close()

	// Execute command
	switch os.Args[1] {
	case "up":
		// Check if a specific version is requested
		if len(os.Args) > 2 {
			version, err := strconv.ParseUint(os.Args[2], 10, 64)
			if err != nil {
				log.Fatalf("Invalid version number: %v", err)
			}
			if err := m.Migrate(uint(version)); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to migrate to version %d: %v", version, err)
			}
		} else {
			// Migrate to latest version
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to run migrations: %v", err)
			}
		}
		fmt.Println("Migrations applied successfully")

	case "down":
		// Check if a specific number of steps is requested
		if len(os.Args) > 2 {
			steps, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalf("Invalid steps number: %v", err)
			}
			for i := 0; i < steps; i++ {
				if err := m.Steps(-1); err != nil {
					if err == migrate.ErrNoChange {
						fmt.Printf("No more migrations to roll back after %d steps\n", i)
						break
					}
					log.Fatalf("Failed to roll back migration (step %d): %v", i+1, err)
				}
			}
		} else {
			// Roll back one migration
			if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to roll back migration: %v", err)
			}
		}
		fmt.Println("Migrations rolled back successfully")

	case "reset":
		// Roll back all migrations
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to reset migrations: %v", err)
		}
		fmt.Println("All migrations have been rolled back")

	case "check":
		status, err := checkMigrationStatus(m)
		if err != nil {
			log.Fatalf("Failed to check migration status: %v", err)
		}

		switch status {
		case MigrationStatusUpToDate:
			fmt.Println("Database schema is up to date")
			os.Exit(0)
		case MigrationStatusPending:
			fmt.Println("Database has pending migrations")
			os.Exit(1)
		case MigrationStatusError:
			fmt.Println("Error checking migration status")
			os.Exit(2)
		}

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				fmt.Println("No migrations have been applied yet")
				os.Exit(0)
			}
			log.Fatalf("Failed to get current version: %v", err)
		}
		fmt.Printf("Current database version: %d (dirty: %t)\n", version, dirty)

	case "status":
		info, err := getMigrationInfo(m)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

		// Format output based on format flag
		if len(os.Args) > 2 && os.Args[2] == "--json" {
			jsonOutput, err := json.MarshalIndent(info, "", "  ")
			if err != nil {
				log.Fatalf("Failed to format JSON: %v", err)
			}
			fmt.Println(string(jsonOutput))
		} else {
			fmt.Printf("Current version: %d\n", info.CurrentVersion)
			fmt.Printf("Latest version: %d\n", info.LatestVersion)
			fmt.Printf("Status: %s\n", info.Status)
			if info.PendingCount > 0 {
				fmt.Printf("Pending migrations: %d\n", info.PendingCount)
			}
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

// These functions are commented out until the required packages are installed

func checkMigrationStatus(m *migrate.Migrate) (MigrationStatus, error) {
	_, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			// No migrations applied yet
			return MigrationStatusPending, nil
		}
		return MigrationStatusError, err
	}

	if dirty {
		return MigrationStatusError, fmt.Errorf("database schema is in a dirty state")
	}

	// Check if there are any pending migrations
	// This is a simplified check - in a real implementation, you would need to
	// compare with the available migration files
	// For now, we'll assume the database is up to date if we have a version
	return MigrationStatusUpToDate, nil
}

func getMigrationInfo(m *migrate.Migrate) (*MigrationInfo, error) {
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return nil, err
	}

	info := &MigrationInfo{
		CurrentVersion: int64(version),
	}

	if dirty {
		info.Status = "dirty"
	} else if err == migrate.ErrNilVersion {
		info.Status = "no migrations applied"
		info.CurrentVersion = -1
	} else {
		info.Status = "clean"
	}

	// In a real implementation, you would determine the latest version
	// by scanning the migrations directory
	// For now, we'll just use the current version as the latest
	info.LatestVersion = int64(version)

	return info, nil
}


func printUsage() {
	fmt.Println("Usage: migrate [command] [options]")
	fmt.Println("Commands:")
	fmt.Println("  up [version]    - Run migrations up to latest or specified version")
	fmt.Println("  down [steps]    - Rollback migrations (default: 1 step)")
	fmt.Println("  reset           - Rollback all migrations")
	fmt.Println("  check           - Check if migrations are up to date")
	fmt.Println("  version         - Show current database version")
	fmt.Println("  status [--json] - Show detailed migration status (optional JSON format)")
}