package handler

import (
	"encoding/json"
	"strconv"

	"gohac/internal/adapter/database"
	"gohac/internal/adapter/repository"
	"gohac/internal/core/domain"
	repoInterface "gohac/internal/core/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PageHandler handles page-related HTTP requests
type PageHandler struct {
	db *gorm.DB
}

// NewPageHandler creates a new page handler instance
func NewPageHandler(db *gorm.DB) *PageHandler {
	return &PageHandler{
		db: db,
	}
}

// CreatePageRequest represents the request body for creating a page
type CreatePageRequest struct {
	Slug   string         `json:"slug" validate:"required"`
	Title  string         `json:"title" validate:"required"`
	Blocks []domain.Block `json:"blocks,omitempty"`
	Status string         `json:"status,omitempty"`
	Meta   map[string]any `json:"meta,omitempty"`
}

// UpdatePageRequest represents the request body for updating a page
type UpdatePageRequest struct {
	Slug   string         `json:"slug,omitempty"`
	Title  string         `json:"title,omitempty"`
	Blocks []domain.Block `json:"blocks,omitempty"`
	Status string         `json:"status,omitempty"`
	Meta   map[string]any `json:"meta,omitempty"`
}

// CreatePage handles POST /api/v1/pages
func (h *PageHandler) CreatePage(c *fiber.Ctx) error {
	var req CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Validate required fields
	if req.Slug == "" || req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Slug and title are required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context (fallback to handler's DB)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		// Fallback to handler's DB for community edition
		db = h.db
	}

	repo := repository.NewPageRepository(db)

	// Get tenant ID from context (empty string for community edition)
	tenantID := ""
	if tenantIDVal := c.Locals("tenant_id"); tenantIDVal != nil {
		if tid, ok := tenantIDVal.(string); ok {
			tenantID = tid
		}
	}

	// Determine status
	status := domain.PageStatusDraft
	if req.Status != "" {
		status = domain.PageStatus(req.Status)
		if status != domain.PageStatusDraft && status != domain.PageStatusPublished && status != domain.PageStatusArchived {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid status. Must be 'draft', 'published', or 'archived'",
				"code":  fiber.StatusBadRequest,
			})
		}
	}

	// Marshal blocks to JSON
	var blocksJSON datatypes.JSON
	if len(req.Blocks) > 0 {
		blocksBytes, err := json.Marshal(req.Blocks)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid blocks format",
				"code":  fiber.StatusBadRequest,
			})
		}
		blocksJSON = blocksBytes
	}

	// Marshal meta to JSON
	var metaJSON datatypes.JSON
	if req.Meta != nil {
		metaBytes, err := json.Marshal(req.Meta)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid meta format",
				"code":  fiber.StatusBadRequest,
			})
		}
		metaJSON = metaBytes
	}

	// Create page
	page := &domain.Page{
		TenantID: tenantID,
		Slug:     req.Slug,
		Title:    req.Title,
		Status:   status,
		Blocks:   blocksJSON,
		Meta:     metaJSON,
	}

	if err := repo.Create(c.Context(), page); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create page",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(page)
}

// ListPages handles GET /api/v1/pages
func (h *PageHandler) ListPages(c *fiber.Ctx) error {
	// Get database from context (fallback to handler's DB)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		// Fallback to handler's DB for community edition
		db = h.db
	}

	repo := repository.NewPageRepository(db)

	// Parse query parameters
	limit := 20 // default
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

	status := c.Query("status")
	search := c.Query("search")

	// Get tenant ID from context
	var tenantID *uuid.UUID
	if tenantIDVal := c.Locals("tenant_id"); tenantIDVal != nil {
		if tenantIDStr, ok := tenantIDVal.(string); ok && tenantIDStr != "" {
			if parsedUUID, err := uuid.Parse(tenantIDStr); err == nil {
				tenantID = &parsedUUID
			}
		}
	}

	opts := repoInterface.ListPageOptions{
		Limit:    limit,
		Offset:   offset,
		Status:   status,
		Search:   search,
		TenantID: tenantID,
	}

	pages, total, err := repo.List(c.Context(), opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list pages",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"data":   pages,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetPage handles GET /api/v1/pages/:id
func (h *PageHandler) GetPage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page ID",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context (fallback to handler's DB)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		// Fallback to handler's DB for community edition
		db = h.db
	}

	repo := repository.NewPageRepository(db)

	page, err := repo.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "page not found: record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Page not found",
				"code":  fiber.StatusNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get page",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(page)
}

// UpdatePage handles PUT /api/v1/pages/:id
func (h *PageHandler) UpdatePage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page ID",
			"code":  fiber.StatusBadRequest,
		})
	}

	var req UpdatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context (fallback to handler's DB)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		// Fallback to handler's DB for community edition
		db = h.db
	}

	repo := repository.NewPageRepository(db)

	// Get existing page
	page, err := repo.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "page not found: record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Page not found",
				"code":  fiber.StatusNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get page",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Update fields if provided
	if req.Slug != "" {
		page.Slug = req.Slug
	}
	if req.Title != "" {
		page.Title = req.Title
	}
	if req.Status != "" {
		status := domain.PageStatus(req.Status)
		if status != domain.PageStatusDraft && status != domain.PageStatusPublished && status != domain.PageStatusArchived {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid status. Must be 'draft', 'published', or 'archived'",
				"code":  fiber.StatusBadRequest,
			})
		}
		page.Status = status
	}
	if req.Blocks != nil {
		blocksJSON, err := json.Marshal(req.Blocks)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid blocks format",
				"code":  fiber.StatusBadRequest,
			})
		}
		page.Blocks = blocksJSON
	}
	if req.Meta != nil {
		metaJSON, err := json.Marshal(req.Meta)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid meta format",
				"code":  fiber.StatusBadRequest,
			})
		}
		page.Meta = metaJSON
	}

	if err := repo.Update(c.Context(), page); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update page",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(page)
}

// DeletePage handles DELETE /api/v1/pages/:id
func (h *PageHandler) DeletePage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page ID",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context (fallback to handler's DB)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		// Fallback to handler's DB for community edition
		db = h.db
	}

	repo := repository.NewPageRepository(db)

	// Check if page exists
	_, err = repo.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "page not found: record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Page not found",
				"code":  fiber.StatusNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get page",
			"code":  fiber.StatusInternalServerError,
		})
	}

	if err := repo.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete page",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
