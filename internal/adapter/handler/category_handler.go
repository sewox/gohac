package handler

import (
	"log"
	"strconv"
	"strings"

	"gohac/internal/adapter/database"
	"gohac/internal/adapter/repository"
	"gohac/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	db *gorm.DB
}

// NewCategoryHandler creates a new category handler instance
func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{
		db: db,
	}
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// CreateCategory handles POST /api/v1/categories (protected endpoint)
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Validate required fields
	if req.Name == "" || req.Slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and slug are required",
			"code":  fiber.StatusBadRequest,
		})
	}

	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	categoryRepo := repository.NewCategoryRepository(db)

	// Check if slug already exists
	_, err = categoryRepo.GetBySlug(c.Context(), req.Slug)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Category with this slug already exists",
			"code":  fiber.StatusConflict,
		})
	}

	category := &domain.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := categoryRepo.Create(c.Context(), category); err != nil {
		log.Printf("Error creating category: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create category",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// ListCategories handles GET /api/v1/categories (protected endpoint)
func (h *CategoryHandler) ListCategories(c *fiber.Ctx) error {
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	categoryRepo := repository.NewCategoryRepository(db)

	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	categories, total, err := categoryRepo.List(c.Context(), limit, offset)
	if err != nil {
		log.Printf("Error listing categories: %v", err)
		// If table doesn't exist, return empty list instead of error
		if strings.Contains(err.Error(), "no such table") || strings.Contains(err.Error(), "doesn't exist") {
			return c.JSON(fiber.Map{
				"data":   []*domain.Category{},
				"total":  0,
				"limit":  limit,
				"offset": offset,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list categories",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"data":   categories,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetCategory handles GET /api/v1/categories/:id (protected endpoint)
func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID format",
			"code":  fiber.StatusBadRequest,
		})
	}

	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	categoryRepo := repository.NewCategoryRepository(db)
	category, err := categoryRepo.GetByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "category not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
				"code":  fiber.StatusNotFound,
			})
		}
		log.Printf("Error getting category: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get category",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(category)
}

// UpdateCategory handles PUT /api/v1/categories/:id (protected endpoint)
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID format",
			"code":  fiber.StatusBadRequest,
		})
	}

	var req UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	categoryRepo := repository.NewCategoryRepository(db)
	category, err := categoryRepo.GetByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "category not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
				"code":  fiber.StatusNotFound,
			})
		}
		log.Printf("Error getting category for update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get category",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Update fields
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		// Check if new slug already exists for another category
		existingCategory, err := categoryRepo.GetBySlug(c.Context(), req.Slug)
		if err == nil && existingCategory.ID != category.ID {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Category with this slug already exists",
				"code":  fiber.StatusConflict,
			})
		}
		category.Slug = req.Slug
	}
	if req.Description != "" || req.Description == "" {
		category.Description = req.Description
	}

	if err := categoryRepo.Update(c.Context(), category); err != nil {
		log.Printf("Error updating category: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update category",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(category)
}

// DeleteCategory handles DELETE /api/v1/categories/:id (protected endpoint)
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID format",
			"code":  fiber.StatusBadRequest,
		})
	}

	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	categoryRepo := repository.NewCategoryRepository(db)
	if err := categoryRepo.Delete(c.Context(), id); err != nil {
		log.Printf("Error deleting category: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete category",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
