//go:build enterprise

package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// connectPostgres creates a PostgreSQL database connection
// For enterprise edition with multi-tenancy support
func connectPostgres() (*gorm.DB, error) {
	// Get connection string from environment
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback to individual environment variables
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

	return db, nil
}
