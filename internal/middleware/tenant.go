package middleware

import (
	"strings"

	"gohac/config"
	"gohac/internal/adapter/database"

	"github.com/gofiber/fiber/v2"
)

// TenantMiddleware extracts tenant information from subdomain or header
// and sets the appropriate database connection in context
func TenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip tenant resolution in community edition
		if !config.SupportsMultiTenancy() {
			// Use default database connection
			return c.Next()
		}

		// Extract tenant ID from subdomain or header
		tenantID := extractTenantID(c)

		if tenantID == "" {
			// Default tenant or error - for now, use empty string
			tenantID = "default"
		}

		// Get tenant-specific database connection
		// Note: In production, you might want to cache these connections
		db, err := database.ConnectForTenant(tenantID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to connect to tenant database",
			})
		}

		// Set database in context
		ctx := database.SetDBInContext(c.Context(), db)
		c.SetUserContext(ctx)

		// Store tenant ID in locals for easy access
		c.Locals("tenant_id", tenantID)

		return c.Next()
	}
}

// extractTenantID extracts tenant ID from request
// Priority: X-Tenant-ID header > subdomain > default
func extractTenantID(c *fiber.Ctx) string {
	// Check header first
	if tenantID := c.Get("X-Tenant-ID"); tenantID != "" {
		return tenantID
	}

	// Extract from subdomain
	host := c.Hostname()
	if host != "" {
		parts := strings.Split(host, ".")
		if len(parts) > 2 {
			// e.g., client1.cms.com -> client1
			return parts[0]
		}
	}

	return ""
}
