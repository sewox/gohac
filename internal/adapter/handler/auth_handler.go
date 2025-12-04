package handler

import (
	"gohac/internal/adapter/database"
	"gohac/internal/adapter/repository"
	"gohac/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler creates a new auth handler instance
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		db: db,
	}
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

// Login handles user authentication with real database lookup and bcrypt
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context (fallback to handler's DB if needed)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		// Fallback to handler's DB
		db = h.db
	}

	// Get user by email
	userRepo := repository.NewUserRepository(db)
	user, err := userRepo.GetByEmail(c.Context(), req.Email)
	if err != nil {
		// Don't reveal if user exists or not (security best practice)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Generate JWT token (valid for 24 hours)
	token, err := middleware.GenerateToken(user.ID.String(), user.Email, 24)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Create HTTP-only cookie (Secure: false for localhost dev, SameSite: Lax)
	cookie := &fiber.Cookie{
		Name:     middleware.AuthTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   24 * 60 * 60, // 24 hours in seconds
		HTTPOnly: true,
		Secure:   false, // Set to false for localhost development
		SameSite: "Lax", // Changed to Lax for better compatibility
	}

	// Set cookie
	c.Cookie(cookie)

	// Return success response
	response := LoginResponse{
		Success: true,
		Message: "Login successful",
	}
	response.User.ID = user.ID.String()
	response.User.Name = user.Name
	response.User.Email = user.Email
	response.User.Role = string(user.Role)

	return c.Status(fiber.StatusOK).JSON(response)
}

// Me returns the current authenticated user's information
// This endpoint is protected and requires valid JWT token
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Get database from context (fallback to handler's DB if needed)
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		// Fallback to handler's DB
		db = h.db
	}

	// Get user from database
	userRepo := repository.NewUserRepository(db)
	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user ID format",
			"code":  fiber.StatusInternalServerError,
		})
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user ID format",
			"code":  fiber.StatusInternalServerError,
		})
	}

	user, err := userRepo.GetByID(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
			"code":  fiber.StatusNotFound,
		})
	}

	// Return user info (excluding password)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":    user.ID.String(),
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Logout handles user logout by clearing the authentication cookie
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Clear the authentication cookie by setting it to expire immediately
	cookie := &fiber.Cookie{
		Name:     middleware.AuthTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Expire immediately
		HTTPOnly: true,
		Secure:   false, // Set to false for localhost development
		SameSite: "Lax",
	}

	// Set cookie to clear it
	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}

// UpdateProfileRequest represents the request body for updating user profile
type UpdateProfileRequest struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"` // Optional - only update if provided
}

// UpdateProfile handles PUT /api/auth/profile (protected endpoint)
// Allows users to update their own profile (name and password only)
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
			"code":  fiber.StatusUnauthorized,
		})
	}

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get database from context (fallback to handler's DB if needed)
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

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user ID format",
			"code":  fiber.StatusInternalServerError,
		})
	}

	user, err := userRepo.GetByID(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
			"code":  fiber.StatusNotFound,
		})
	}

	// Update fields (only name and password allowed for profile update)
	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Password != "" {
		// Hash new password
		user.Password = req.Password
		if err := user.HashPassword(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process password",
				"code":  fiber.StatusInternalServerError,
			})
		}
	}

	// Save updated user
	if err := userRepo.Update(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Return updated user info (excluding password)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
		"user": fiber.Map{
			"id":    user.ID.String(),
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
