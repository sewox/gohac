//go:build community || !enterprise

package database

import (
	"fmt"

	"gorm.io/gorm"
)

// ConnectForTenant is a stub for community edition
// Multi-tenancy is not supported in community edition
func ConnectForTenant(tenantID string) (*gorm.DB, error) {
	return nil, fmt.Errorf("multi-tenancy is not supported in community edition")
}
