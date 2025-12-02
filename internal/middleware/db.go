package middleware

import (
	"gohac/config"
	"gohac/internal/adapter/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DBMiddleware sets the database connection in context for community edition
// In enterprise mode, TenantMiddleware handles this
func DBMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only set DB in context for community edition
		// Enterprise edition uses TenantMiddleware
		if !config.SupportsMultiTenancy() {
			ctx := database.SetDBInContext(c.Context(), db)
			c.SetUserContext(ctx)
		}
		return c.Next()
	}
}
