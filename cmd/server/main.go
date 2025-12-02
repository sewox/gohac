package main

import (
	"log"
	"os"

	"gohac/config"
	"gohac/internal/adapter/database"
	"gohac/internal/adapter/handler"
	"gohac/internal/core/domain"
	"gohac/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func main() {
	// Initialize database connection
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	if err := migrateDatabase(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Gohac CMS",
		ServerHeader: "Gohac",
		ErrorHandler: errorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			// In development, allow all origins
			// In production, validate against allowed origins list
			if os.Getenv("ENV") == "production" {
				allowedOrigins := []string{
					"localhost:3131",
					// Add your production domains here
				}
				for _, allowed := range allowedOrigins {
					if origin == allowed {
						return true
					}
				}
				return false
			}
			// Development: allow all origins
			return true
		},
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Tenant-ID",
		AllowCredentials: true,
	}))

	// Tenant middleware (only in enterprise mode)
	if config.SupportsMultiTenancy() {
		app.Use(middleware.TenantMiddleware())
	}

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":       "ok",
			"version":      "1.0.0",
			"edition":      getEdition(),
			"database":     config.GetDatabaseDriver(),
			"multi_tenant": config.SupportsMultiTenancy(),
		})
	})

	// Placeholder CMS route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Gohac CMS",
			"edition": getEdition(),
			"docs":    "/api/docs",
		})
	})

	// Setup API routes
	setupAPIRoutes(app, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3131"
	}

	log.Printf("🚀 Gohac CMS Server starting on port %s", port)
	log.Printf("📦 Edition: %s", getEdition())
	log.Printf("💾 Database: %s", config.GetDatabaseDriver())
	log.Printf("🏢 Multi-tenancy: %v", config.SupportsMultiTenancy())

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// migrateDatabase runs database migrations
func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Page{},
		// Add more models here as they are created
	)
}

// setupAPIRoutes sets up API route handlers
func setupAPIRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	// Public auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	// Protected routes (require authentication)
	v1 := api.Group("/v1")
	v1.Use(middleware.Protected()) // Apply auth middleware to all v1 routes

	// Pages endpoint (protected)
	v1.Get("/pages", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		tenantID := c.Locals("tenant_id")

		return c.JSON(fiber.Map{
			"message":   "Pages API endpoint",
			"user_id":   userID,
			"tenant_id": tenantID,
		})
	})
}

// errorHandler is the global error handler
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error": message,
		"code":  code,
	})
}

// getEdition returns the current edition name
func getEdition() string {
	if config.IsEnterprise() {
		return "Enterprise"
	}
	return "Community"
}
