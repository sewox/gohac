package main

import (
	"log"
	"os"

	"gohac/config"
	"gohac/internal/adapter/database"
	"gohac/internal/adapter/handler"
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

	// Run database migrations
	if err := database.Migrate(db); err != nil {
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

	// Static file serving for uploaded files
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}
	// Serve static files from storage directory
	// This allows accessing uploaded files via /static/uploads/images/filename.jpg
	app.Static("/static", storagePath)
	// Also serve from /uploads for convenience (redirects to /static/uploads)
	app.Static("/uploads", storagePath+"/uploads")

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

// setupAPIRoutes sets up API route handlers
func setupAPIRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	// Public auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	// Protected auth routes
	authProtected := api.Group("/auth")
	authProtected.Use(middleware.Protected())
	authProtected.Get("/me", handler.Me)

	// Protected routes (require authentication)
	v1 := api.Group("/v1")
	v1.Use(middleware.Protected()) // Apply auth middleware to all v1 routes

	// Set DB in context for community edition (enterprise uses TenantMiddleware)
	if !config.SupportsMultiTenancy() {
		v1.Use(middleware.DBMiddleware(db))
	}

	// Create page handler
	pageHandler := handler.NewPageHandler(db)

	// Page routes
	v1.Post("/pages", pageHandler.CreatePage)
	v1.Get("/pages", pageHandler.ListPages)
	v1.Get("/pages/:id", pageHandler.GetPage)
	v1.Put("/pages/:id", pageHandler.UpdatePage)
	v1.Delete("/pages/:id", pageHandler.DeletePage)

	// Upload handler
	uploadHandler := handler.NewUploadHandler()
	v1.Post("/upload", uploadHandler.UploadFile)
	v1.Post("/upload/from-url", uploadHandler.DownloadFromURL)

	// Settings handler
	settingsHandler := handler.NewSettingsHandler(db)
	v1.Put("/settings", settingsHandler.UpdateSettings)

	// Menu handler
	menuHandler := handler.NewMenuHandler(db)
	v1.Post("/menus", menuHandler.CreateMenu)
	v1.Get("/menus", menuHandler.ListMenus)
	v1.Get("/menus/:id", menuHandler.GetMenu)
	v1.Put("/menus/:id", menuHandler.UpdateMenu)
	v1.Delete("/menus/:id", menuHandler.DeleteMenu)

	// Public API routes (no authentication required)
	public := app.Group("/api/public")
	// Set DB in context for public routes (community edition)
	if !config.SupportsMultiTenancy() {
		public.Use(middleware.DBMiddleware(db))
	}
	// Register specific routes BEFORE wildcard routes to avoid route conflicts
	public.Get("/settings", settingsHandler.GetSettings)
	public.Get("/menus/:id", menuHandler.GetMenu) // Public endpoint to get menu by ID
	public.Get("/pages/*", pageHandler.GetPageBySlugPublic)
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
