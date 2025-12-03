package handler

import (
	"gohac/internal/adapter/database"
	"gohac/internal/adapter/repository"
	"gohac/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SettingsHandler handles settings and menu-related HTTP requests
type SettingsHandler struct {
	db *gorm.DB
}

// NewSettingsHandler creates a new settings handler instance
func NewSettingsHandler(db *gorm.DB) *SettingsHandler {
	return &SettingsHandler{
		db: db,
	}
}

// GetSettings handles GET /api/public/settings (public endpoint)
func (h *SettingsHandler) GetSettings(c *fiber.Ctx) error {
	// Get database from context (fallback to handler's DB)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	repo := repository.NewSettingsRepository(db)
	settings, err := repo.GetGlobalSettings(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get settings",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(settings)
}

// UpdateSettingsRequest represents the request body for updating settings
type UpdateSettingsRequest struct {
	SiteName     string `json:"site_name"`
	Logo         string `json:"logo"`
	Favicon      string `json:"favicon"`
	ContactEmail string `json:"contact_email"`
	HeaderMenuID string `json:"header_menu_id,omitempty"` // UUID of the menu to display in header
	FooterMenuID string `json:"footer_menu_id,omitempty"` // UUID of the menu to display in footer
}

// UpdateSettings handles PUT /api/v1/settings (protected endpoint)
func (h *SettingsHandler) UpdateSettings(c *fiber.Ctx) error {
	var req UpdateSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context (fallback to handler's DB)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	repo := repository.NewSettingsRepository(db)

	settings := &domain.GlobalSettings{
		SiteName:     req.SiteName,
		Logo:         req.Logo,
		Favicon:      req.Favicon,
		ContactEmail: req.ContactEmail,
		HeaderMenuID: req.HeaderMenuID,
		FooterMenuID: req.FooterMenuID,
	}

	if err := repo.UpdateGlobalSettings(c.Context(), settings); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update settings",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(settings)
}
