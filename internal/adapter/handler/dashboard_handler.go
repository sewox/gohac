package handler

import (
	"log"
	"os"
	"path/filepath"

	"gohac/internal/adapter/database"
	"gohac/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DashboardHandler handles dashboard-related HTTP requests
type DashboardHandler struct {
	db *gorm.DB
}

// NewDashboardHandler creates a new dashboard handler instance
func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{
		db: db,
	}
}

// GetStats handles GET /api/v1/dashboard/stats (protected endpoint)
// Returns statistics about pages, users, and media files
func (h *DashboardHandler) GetStats(c *fiber.Ctx) error {
	// Get database from context (fallback to handler's DB if needed)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	// Count pages
	var pageCount int64
	if err := db.Model(&domain.Page{}).Count(&pageCount).Error; err != nil {
		log.Printf("Error counting pages: %v", err)
		pageCount = 0
	}

	// Count users
	var userCount int64
	if err := db.Model(&domain.User{}).Count(&userCount).Error; err != nil {
		log.Printf("Error counting users: %v", err)
		userCount = 0
	}

	// Count media files
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}
	uploadPath := filepath.Join(storagePath, "uploads")

	var mediaCount int64 = 0
	// Check if upload directory exists
	if _, err := os.Stat(uploadPath); err == nil {
		if entries, err := os.ReadDir(uploadPath); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					mediaCount++
				}
			}
		}
	}

	// Count posts (if blog system is enabled)
	var postCount int64 = 0
	if db.Migrator().HasTable(&domain.Post{}) {
		if err := db.Model(&domain.Post{}).Count(&postCount).Error; err != nil {
			log.Printf("Error counting posts: %v", err)
			postCount = 0
		}
	}

	// Count categories (if blog system is enabled)
	var categoryCount int64 = 0
	if db.Migrator().HasTable(&domain.Category{}) {
		if err := db.Model(&domain.Category{}).Count(&categoryCount).Error; err != nil {
			log.Printf("Error counting categories: %v", err)
			categoryCount = 0
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"stats": fiber.Map{
			"pages":      pageCount,
			"users":      userCount,
			"media":      mediaCount,
			"posts":      postCount,
			"categories": categoryCount,
		},
	})
}
