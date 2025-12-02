//go:build enterprise

package database

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectForTenant creates a database connection for a specific tenant
// In enterprise mode, this can switch between databases or schemas
func ConnectForTenant(tenantID string) (*gorm.DB, error) {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "postgres" // Default for enterprise
	}

	switch driver {
	case "sqlite":
		return connectSQLiteForTenant(tenantID)
	case "postgres":
		return connectPostgresForTenant(tenantID)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}

// connectSQLiteForTenant creates a tenant-specific SQLite database
func connectSQLiteForTenant(tenantID string) (*gorm.DB, error) {
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, fmt.Sprintf("%s.db", tenantID))

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite for tenant %s: %w", tenantID, err)
	}

	return db, nil
}

// connectPostgresForTenant creates a tenant-specific PostgreSQL connection
// Uses schema-based multi-tenancy (each tenant has its own schema)
func connectPostgresForTenant(tenantID string) (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := getEnvOrDefault("DB_HOST", "localhost")
		port := getEnvOrDefault("DB_PORT", "5432")
		user := getEnvOrDefault("DB_USER", "postgres")
		password := getEnvOrDefault("DB_PASSWORD", "")
		dbname := getEnvOrDefault("DB_NAME", "gohac")
		sslmode := getEnvOrDefault("DB_SSLMODE", "disable")

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Set search_path to tenant-specific schema
	schemaName := fmt.Sprintf("tenant_%s", tenantID)
	if err := db.Exec(fmt.Sprintf("SET search_path TO %s", schemaName)).Error; err != nil {
		// Schema might not exist, create it
		if err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)).Error; err != nil {
			return nil, fmt.Errorf("failed to create/switch to schema %s: %w", schemaName, err)
		}
	}

	return db, nil
}
