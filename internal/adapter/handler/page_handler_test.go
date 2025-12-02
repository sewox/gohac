package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"gohac/internal/adapter/database"
	"gohac/internal/adapter/repository"
	"gohac/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestApp creates a Fiber app with test database and handlers
func setupTestApp(t *testing.T) (*fiber.App, *gorm.DB) {
	// Create in-memory test database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate schema
	err = db.AutoMigrate(&domain.Page{})
	require.NoError(t, err)

	// Create Fiber app
	app := fiber.New()

	// Create page handler
	pageHandler := NewPageHandler(db)

	// Setup routes
	v1 := app.Group("/api/v1")
	v1.Post("/pages", pageHandler.CreatePage)
	v1.Get("/pages", pageHandler.ListPages)
	v1.Get("/pages/:id", pageHandler.GetPage)
	v1.Put("/pages/:id", pageHandler.UpdatePage)
	v1.Delete("/pages/:id", pageHandler.DeletePage)

	return app, db
}

func TestPageHandler_CreatePage_ValidData(t *testing.T) {
	app, _ := setupTestApp(t)

	reqBody := CreatePageRequest{
		Slug:   "test-page",
		Title:  "Test Page",
		Status: "draft",
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/pages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var page domain.Page
	err = json.NewDecoder(resp.Body).Decode(&page)
	require.NoError(t, err)
	assert.Equal(t, "test-page", page.Slug)
	assert.Equal(t, "Test Page", page.Title)
	assert.Equal(t, domain.PageStatusDraft, page.Status)
	assert.NotEqual(t, uuid.Nil, page.ID)
}

func TestPageHandler_CreatePage_InvalidData(t *testing.T) {
	app, _ := setupTestApp(t)

	// Test with missing required fields
	reqBody := map[string]interface{}{
		"slug": "test-page",
		// Missing title
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/pages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var errorResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Contains(t, errorResp["error"], "required")
}

func TestPageHandler_CreatePage_InvalidStatus(t *testing.T) {
	app, _ := setupTestApp(t)

	reqBody := CreatePageRequest{
		Slug:   "test-page",
		Title:  "Test Page",
		Status: "invalid-status",
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/pages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestPageHandler_ListPages(t *testing.T) {
	app, db := setupTestApp(t)

	// Create some test pages
	repo := repository.NewPageRepository(db)
	ctx := database.SetDBInContext(context.Background(), db)

	page1 := &domain.Page{
		Slug:   "page-1",
		Title:  "Page 1",
		Status: domain.PageStatusDraft,
	}
	err := repo.Create(ctx, page1)
	require.NoError(t, err)

	page2 := &domain.Page{
		Slug:   "page-2",
		Title:  "Page 2",
		Status: domain.PageStatusPublished,
	}
	err = repo.Create(ctx, page2)
	require.NoError(t, err)

	// Test list endpoint
	req := httptest.NewRequest("GET", "/api/v1/pages", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "total")
	assert.Equal(t, float64(2), response["total"])
}

func TestPageHandler_GetPage_NotFound(t *testing.T) {
	app, _ := setupTestApp(t)

	nonExistentID := uuid.New()
	req := httptest.NewRequest("GET", "/api/v1/pages/"+nonExistentID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var errorResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Contains(t, errorResp["error"], "not found")
}

func TestPageHandler_GetPage_InvalidID(t *testing.T) {
	app, _ := setupTestApp(t)

	req := httptest.NewRequest("GET", "/api/v1/pages/invalid-uuid", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestPageHandler_UpdatePage(t *testing.T) {
	app, db := setupTestApp(t)

	// Create a page first
	repo := repository.NewPageRepository(db)
	ctx := database.SetDBInContext(context.Background(), db)

	page := &domain.Page{
		Slug:   "original-slug",
		Title:  "Original Title",
		Status: domain.PageStatusDraft,
	}
	err := repo.Create(ctx, page)
	require.NoError(t, err)

	// Update the page
	updateReq := UpdatePageRequest{
		Title:  "Updated Title",
		Status: "published",
	}

	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "/api/v1/pages/"+page.ID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var updatedPage domain.Page
	err = json.NewDecoder(resp.Body).Decode(&updatedPage)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedPage.Title)
	assert.Equal(t, domain.PageStatusPublished, updatedPage.Status)
}

func TestPageHandler_DeletePage(t *testing.T) {
	app, db := setupTestApp(t)

	// Create a page first
	repo := repository.NewPageRepository(db)
	ctx := database.SetDBInContext(context.Background(), db)

	page := &domain.Page{
		Slug:   "to-delete",
		Title:  "Page to Delete",
		Status: domain.PageStatusDraft,
	}
	err := repo.Create(ctx, page)
	require.NoError(t, err)

	// Delete the page
	req := httptest.NewRequest("DELETE", "/api/v1/pages/"+page.ID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	// Verify page is deleted
	_, err = repo.GetByID(ctx, page.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
