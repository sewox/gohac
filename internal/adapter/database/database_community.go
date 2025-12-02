//go:build community || !enterprise

package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// DBKey is the context key for storing database instance
type DBKey string

const DBContextKey DBKey = "db"

// Connect creates a database connection based on build tags
func Connect() (*gorm.DB, error) {
	return connectSQLite()
}

// GetDBFromContext retrieves the database instance from context
func GetDBFromContext(ctx context.Context) (*gorm.DB, error) {
	db, ok := ctx.Value(DBContextKey).(*gorm.DB)
	if !ok || db == nil {
		return nil, fmt.Errorf("database not found in context")
	}
	return db, nil
}

// SetDBInContext sets the database instance in context
func SetDBInContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, DBContextKey, db)
}
