package handler

import (
	"encoding/json"
	"log"
	"strconv"

	"gohac/internal/adapter/database"
	"gohac/internal/adapter/repository"
	"gohac/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// MenuHandler handles menu-related HTTP requests
type MenuHandler struct {
	db *gorm.DB
}

// NewMenuHandler creates a new menu handler instance
func NewMenuHandler(db *gorm.DB) *MenuHandler {
	return &MenuHandler{
		db: db,
	}
}

// CreateMenuRequest represents the request body for creating a menu
type CreateMenuRequest struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description,omitempty"`
	Items       []domain.MenuItem `json:"items,omitempty"`
}

// UpdateMenuRequest represents the request body for updating a menu
type UpdateMenuRequest struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Items       []domain.MenuItem `json:"items,omitempty"`
}

// CreateMenu handles POST /api/v1/menus (protected endpoint)
func (h *MenuHandler) CreateMenu(c *fiber.Ctx) error {
	var req CreateMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context (fallback to handler's DB)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	repo := repository.NewMenuRepository(db)

	// Get tenant ID from context (empty string for community edition)
	tenantID := ""
	if tenantIDVal := c.Locals("tenant_id"); tenantIDVal != nil {
		if tid, ok := tenantIDVal.(string); ok {
			tenantID = tid
		}
	}

	// Marshal menu items to JSON
	itemsJSON := []byte("[]")
	if len(req.Items) > 0 {
		itemsBytes, err := json.Marshal(req.Items)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid menu items format",
				"code":  fiber.StatusBadRequest,
			})
		}
		itemsJSON = itemsBytes
	}

	menu := &domain.Menu{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		Items:       datatypes.JSON(itemsJSON),
	}

	if err := repo.Create(c.Context(), menu); err != nil {
		// Log the actual error for debugging
		log.Printf("Error creating menu: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create menu: " + err.Error(),
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Parse items for response
	var items []domain.MenuItem
	if len(menu.Items) > 0 {
		json.Unmarshal(menu.Items, &items)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":          menu.ID,
		"name":        menu.Name,
		"description": menu.Description,
		"items":       items,
		"created_at":  menu.CreatedAt,
		"updated_at":  menu.UpdatedAt,
	})
}

// ListMenus handles GET /api/v1/menus (protected endpoint)
func (h *MenuHandler) ListMenus(c *fiber.Ctx) error {
	// Get database from context (fallback to handler's DB)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	repo := repository.NewMenuRepository(db)

	// Parse query parameters
	limit := 50 // default
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	menus, total, err := repo.List(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list menus",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Parse items for each menu
	result := make([]fiber.Map, len(menus))
	for i, menu := range menus {
		var items []domain.MenuItem
		if len(menu.Items) > 0 {
			json.Unmarshal(menu.Items, &items)
		}
		result[i] = fiber.Map{
			"id":          menu.ID,
			"name":        menu.Name,
			"description": menu.Description,
			"items":       items,
			"created_at":  menu.CreatedAt,
			"updated_at":  menu.UpdatedAt,
		}
	}

	return c.JSON(fiber.Map{
		"data":   result,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetMenu handles GET /api/v1/menus/:id (protected endpoint) and GET /api/public/menus/:id (public endpoint)
func (h *MenuHandler) GetMenu(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid menu ID",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context (fallback to handler's DB)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	repo := repository.NewMenuRepository(db)

	menu, err := repo.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "menu not found: record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Menu not found",
				"code":  fiber.StatusNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get menu",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Parse items for response
	var items []domain.MenuItem
	if len(menu.Items) > 0 {
		json.Unmarshal(menu.Items, &items)
	}

	return c.JSON(fiber.Map{
		"id":          menu.ID,
		"name":        menu.Name,
		"description": menu.Description,
		"items":       items,
		"created_at":  menu.CreatedAt,
		"updated_at":  menu.UpdatedAt,
	})
}

// UpdateMenu handles PUT /api/v1/menus/:id (protected endpoint)
func (h *MenuHandler) UpdateMenu(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid menu ID",
			"code":  fiber.StatusBadRequest,
		})
	}

	var req UpdateMenuRequest
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

	repo := repository.NewMenuRepository(db)

	// Get existing menu
	menu, err := repo.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "menu not found: record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Menu not found",
				"code":  fiber.StatusNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get menu",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Update fields
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Description != "" || c.Query("clear_description") == "true" {
		menu.Description = req.Description
	}
	if req.Items != nil {
		itemsBytes, err := json.Marshal(req.Items)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid menu items format",
				"code":  fiber.StatusBadRequest,
			})
		}
		menu.Items = datatypes.JSON(itemsBytes)
	}

	if err := repo.Update(c.Context(), menu); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update menu",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Parse items for response
	var items []domain.MenuItem
	if len(menu.Items) > 0 {
		json.Unmarshal(menu.Items, &items)
	}

	return c.JSON(fiber.Map{
		"id":          menu.ID,
		"name":        menu.Name,
		"description": menu.Description,
		"items":       items,
		"created_at":  menu.CreatedAt,
		"updated_at":  menu.UpdatedAt,
	})
}

// DeleteMenu handles DELETE /api/v1/menus/:id (protected endpoint)
func (h *MenuHandler) DeleteMenu(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid menu ID",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context (fallback to handler's DB)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	repo := repository.NewMenuRepository(db)

	// Check if menu exists
	_, err = repo.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "menu not found: record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Menu not found",
				"code":  fiber.StatusNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get menu",
			"code":  fiber.StatusInternalServerError,
		})
	}

	if err := repo.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete menu",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
