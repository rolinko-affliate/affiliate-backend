# Database Migration Command

This module provides a command-line tool for managing database migrations. It uses the golang-migrate library to apply, rollback, and check the status of migrations.

## Key Components

### Main Function

The `main` function is the entry point of the migration tool:

```go
func main() {
    // Load configuration
    config.LoadConfig()
    dbURL := config.AppConfig.DatabaseURL

    // Parse command-line arguments
    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }

    command := os.Args[1]

    // Execute the appropriate migration command
    switch command {
    case "up":
        migrateUp(dbURL)
    case "down":
        migrateDown(dbURL)
    case "reset":
        migrateReset(dbURL)
    case "version":
        migrateVersion(dbURL)
    case "status":
        migrateStatus(dbURL)
    case "check":
        checkMigrations(dbURL)
    case "create":
        createMigration(os.Args[2:])
    default:
        fmt.Printf("Unknown command: %s\n", command)
        printUsage()
        os.Exit(1)
    }
}
```

### Migration Commands

The module provides several migration commands:

- `up`: Apply all pending migrations
- `down`: Rollback the most recent migration
- `reset`: Rollback all migrations
- `version`: Show current database version
- `status`: Show detailed migration status
- `check`: Check if migrations are up to date
- `create`: Create a new migration file

### Migration Functions

Each command is implemented as a function:

```go
func migrateUp(dbURL string) {
    m, err := getMigrator(dbURL)
    if err != nil {
        log.Fatalf("Error creating migrator: %v", err)
    }
    defer closeMigrator(m)

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Error applying migrations: %v", err)
    }
    fmt.Println("Migrations applied successfully")
}

// Other migration functions...
```

### Utility Functions

The module includes utility functions for working with migrations:

- `getMigrator`: Creates a new migrate instance
- `closeMigrator`: Closes the migrate instance
- `printUsage`: Prints usage information
- `createMigration`: Creates new migration files

## Migration Files

Migration files are stored in the `migrations` directory and follow the naming convention:

```
000001_name.up.sql   # SQL to apply the migration
000001_name.down.sql # SQL to rollback the migration
```

## Usage

The migration tool is used from the command line:

```bash
# Apply all pending migrations
go run cmd/migrate/main.go up

# Rollback the most recent migration
go run cmd/migrate/main.go down

# Rollback all migrations
go run cmd/migrate/main.go reset

# Show current database version
go run cmd/migrate/main.go version

# Show detailed migration status
go run cmd/migrate/main.go status

# Check if migrations are up to date
go run cmd/migrate/main.go check

# Create a new migration file
go run cmd/migrate/main.go create add_new_table
```

## Integration with Makefile

The migration tool is integrated with the project's Makefile:

```makefile
# Apply all pending migrations
migrate-up:
	go run cmd/migrate/main.go up

# Rollback the most recent migration
migrate-down:
	go run cmd/migrate/main.go down

# Other migration commands...
```

## Error Handling

The migration tool provides detailed error messages:

- Connection errors when the database is unavailable
- Migration errors when SQL statements fail
- File system errors when creating migration files
- Version errors when the migration version is dirty