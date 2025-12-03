package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gohac/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Auto-migrate
	db.AutoMigrate(&domain.User{})

	return db
}

func TestUserHandler_CreateUser(t *testing.T) {
	db := setupTestDB()
	handler := NewUserHandler(db)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		// Mock admin user in context
		c.Locals("user_id", uuid.New().String())
		c.Locals("user_email", "admin@test.com")
		// Set DB in context
		c.Locals("db", db)
		return c.Next()
	})

	// Mock requireAdmin to always pass (we'll test admin check separately)
	app.Post("/users", func(c *fiber.Ctx) error {
		// Bypass admin check for this test
		return handler.CreateUser(c)
	})

	// Create a test admin user first (for requireAdmin check)
	adminUser := &domain.User{
		Name:     "Admin",
		Email:    "admin@test.com",
		Password: "admin123",
		Role:     domain.UserRoleAdmin,
	}
	adminUser.HashPassword()
	db.Create(adminUser)

	// Test creating a new user
	reqBody := CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "editor",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Verify user was created
	var user domain.User
	db.Where("email = ?", "test@example.com").First(&user)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, domain.UserRoleEditor, user.Role)
	assert.True(t, user.CheckPassword("password123"))
}

func TestUserHandler_ListUsers(t *testing.T) {
	db := setupTestDB()
	handler := NewUserHandler(db)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	// Create test users
	user1 := &domain.User{
		Name:     "User 1",
		Email:    "user1@test.com",
		Password: "pass123",
		Role:     domain.UserRoleEditor,
	}
	user1.HashPassword()
	db.Create(user1)

	user2 := &domain.User{
		Name:     "User 2",
		Email:    "user2@test.com",
		Password: "pass123",
		Role:     domain.UserRoleAdmin,
	}
	user2.HashPassword()
	db.Create(user2)

	// Create admin user for auth
	adminUser := &domain.User{
		Name:     "Admin",
		Email:    "admin@test.com",
		Password: "admin123",
		Role:     domain.UserRoleAdmin,
	}
	adminUser.HashPassword()
	db.Create(adminUser)

	app.Get("/users", func(c *fiber.Ctx) error {
		c.Locals("user_id", adminUser.ID.String())
		return handler.ListUsers(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(t, result["data"])
}

func TestUserHandler_UpdateUser(t *testing.T) {
	db := setupTestDB()
	handler := NewUserHandler(db)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	// Create admin user
	adminUser := &domain.User{
		Name:     "Admin",
		Email:    "admin@test.com",
		Password: "admin123",
		Role:     domain.UserRoleAdmin,
	}
	adminUser.HashPassword()
	db.Create(adminUser)

	// Create user to update
	user := &domain.User{
		Name:     "Original Name",
		Email:    "original@test.com",
		Password: "pass123",
		Role:     domain.UserRoleEditor,
	}
	user.HashPassword()
	db.Create(user)

	app.Put("/users/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", adminUser.ID.String())
		return handler.UpdateUser(c)
	})

	reqBody := UpdateUserRequest{
		Name:  "Updated Name",
		Email: "updated@test.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/users/"+user.ID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify user was updated
	var updatedUser domain.User
	db.First(&updatedUser, user.ID)
	assert.Equal(t, "Updated Name", updatedUser.Name)
	assert.Equal(t, "updated@test.com", updatedUser.Email)
}

func TestUserHandler_DeleteUser(t *testing.T) {
	db := setupTestDB()
	handler := NewUserHandler(db)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	// Create admin user
	adminUser := &domain.User{
		Name:     "Admin",
		Email:    "admin@test.com",
		Password: "admin123",
		Role:     domain.UserRoleAdmin,
	}
	adminUser.HashPassword()
	db.Create(adminUser)

	// Create user to delete
	user := &domain.User{
		Name:     "To Delete",
		Email:    "delete@test.com",
		Password: "pass123",
		Role:     domain.UserRoleEditor,
	}
	user.HashPassword()
	db.Create(user)

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", adminUser.ID.String())
		return handler.DeleteUser(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/users/"+user.ID.String(), nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify user was deleted
	var deletedUser domain.User
	err = db.First(&deletedUser, user.ID).Error
	assert.Error(t, err) // Should not find the user
}
