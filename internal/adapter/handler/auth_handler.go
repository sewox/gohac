package handler

import (
	"gohac/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

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
		Email string `json:"email"`
	} `json:"user"`
}

// Login handles user authentication
// Mock authentication: accepts any email with password "admin123"
func Login(c *fiber.Ctx) error {
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

	// Mock authentication check
	// TODO: Replace with actual database lookup and password hashing
	if req.Password != "admin123" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Generate JWT token (valid for 24 hours)
	token, err := middleware.GenerateToken(req.Email, req.Email, 24)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Create HTTP-only, Secure, SameSite=Strict cookie
	cookie := &fiber.Cookie{
		Name:     middleware.AuthTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   24 * 60 * 60, // 24 hours in seconds
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	}

	// Set cookie
	c.Cookie(cookie)

	// Return success response
	response := LoginResponse{
		Success: true,
		Message: "Login successful",
	}
	response.User.ID = req.Email
	response.User.Email = req.Email

	return c.Status(fiber.StatusOK).JSON(response)
}
