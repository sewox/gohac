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

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new user handler instance
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

// requireAdmin checks if the current user has admin role
func (h *UserHandler) requireAdmin(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Get database from context
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	// Get current user
	userRepo := repository.NewUserRepository(db)
	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user ID format",
			"code":  fiber.StatusInternalServerError,
		})
	}

	user, err := userRepo.GetByID(c.Context(), userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
			"code":  fiber.StatusUnauthorized,
		})
	}

	if user.Role != domain.UserRoleAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin role required",
			"code":  fiber.StatusForbidden,
		})
	}

	return nil
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin editor"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"` // Optional - only update if provided
	Role     string `json:"role,omitempty"`
}

// ListUsers handles GET /api/v1/users (protected endpoint, admin only)
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	// Check admin role
	if err := h.requireAdmin(c); err != nil {
		return err
	}

	// Get database from context
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	repo := repository.NewUserRepository(db)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	users, total, err := repo.List(c.Context(), limit, offset)
	if err != nil {
		log.Printf("Error listing users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list users",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Prepare response (exclude passwords)
	var responseUsers []fiber.Map
	for _, user := range users {
		responseUsers = append(responseUsers, fiber.Map{
			"id":         user.ID.String(),
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"data":   responseUsers,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetUser handles GET /api/v1/users/:id (protected endpoint, admin only)
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	// Check admin role
	if err := h.requireAdmin(c); err != nil {
		return err
	}

	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("Error parsing user ID for get: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	repo := repository.NewUserRepository(db)
	user, err := repo.GetByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
				"code":  fiber.StatusNotFound,
			})
		}
		log.Printf("Error getting user by ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"id":         user.ID.String(),
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// CreateUser handles POST /api/v1/users (protected endpoint, admin only)
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Check admin role
	if err := h.requireAdmin(c); err != nil {
		return err
	}

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, email, and password are required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Validate role
	role := domain.UserRole(strings.ToLower(req.Role))
	if role != domain.UserRoleAdmin && role != domain.UserRoleEditor {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role. Must be 'admin' or 'editor'",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	repo := repository.NewUserRepository(db)

	// Check if email already exists
	_, err = repo.GetByEmail(c.Context(), req.Email)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User with this email already exists",
			"code":  fiber.StatusConflict,
		})
	}

	// Create user
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // Will be hashed
		Role:     role,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
			"code":  fiber.StatusInternalServerError,
		})
	}

	if err := repo.Create(c.Context(), user); err != nil {
		log.Printf("Error creating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         user.ID.String(),
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// UpdateUser handles PUT /api/v1/users/:id (protected endpoint, admin only)
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Check admin role
	if err := h.requireAdmin(c); err != nil {
		return err
	}

	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
			"code":  fiber.StatusBadRequest,
		})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	repo := repository.NewUserRepository(db)
	user, err := repo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
			"code":  fiber.StatusNotFound,
		})
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if email already exists (excluding current user)
		existingUser, err := repo.GetByEmail(c.Context(), req.Email)
		if err == nil && existingUser.ID != user.ID {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "User with this email already exists",
				"code":  fiber.StatusConflict,
			})
		}
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password = req.Password
		if err := user.HashPassword(); err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update user",
				"code":  fiber.StatusInternalServerError,
			})
		}
	}
	if req.Role != "" {
		role := domain.UserRole(strings.ToLower(req.Role))
		if role != domain.UserRoleAdmin && role != domain.UserRoleEditor {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid role. Must be 'admin' or 'editor'",
				"code":  fiber.StatusBadRequest,
			})
		}
		user.Role = role
	}

	if err := repo.Update(c.Context(), user); err != nil {
		log.Printf("Error updating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"id":         user.ID.String(),
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// DeleteUser handles DELETE /api/v1/users/:id (protected endpoint, admin only)
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	// Check admin role
	if err := h.requireAdmin(c); err != nil {
		return err
	}

	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Prevent deleting yourself
	userID := c.Locals("user_id")
	if userID != nil {
		userIDStr, ok := userID.(string)
		if ok && userIDStr == id.String() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot delete your own account",
				"code":  fiber.StatusBadRequest,
			})
		}
	}

	// Get database from context
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	repo := repository.NewUserRepository(db)
	if err := repo.Delete(c.Context(), id); err != nil {
		log.Printf("Error deleting user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
