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

func TestPageHandler_UpdatePage_WithBlocks(t *testing.T) {
	app, db := setupTestApp(t)

	// Create a page first
	repo := repository.NewPageRepository(db)
	ctx := database.SetDBInContext(context.Background(), db)

	page := &domain.Page{
		Slug:   "blocks-test",
		Title:  "Blocks Test Page",
		Status: domain.PageStatusDraft,
	}
	err := repo.Create(ctx, page)
	require.NoError(t, err)

	// Update the page with blocks
	blocks := []domain.Block{
		{
			ID:   "block-1",
			Type: "hero",
			Data: json.RawMessage(`{"title":"Hello","subtitle":"World"}`),
		},
	}

	updateReq := UpdatePageRequest{
		Title:  "Updated Title",
		Blocks: blocks,
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

	// Verify blocks were saved
	var savedBlocks []domain.Block
	err = json.Unmarshal(updatedPage.Blocks, &savedBlocks)
	require.NoError(t, err)
	assert.Len(t, savedBlocks, 1)
	assert.Equal(t, "hero", savedBlocks[0].Type)

	// Verify block data
	var heroData domain.HeroBlockData
	err = json.Unmarshal(savedBlocks[0].Data, &heroData)
	require.NoError(t, err)
	assert.Equal(t, "Hello", heroData.Title)
	assert.Equal(t, "World", heroData.Subtitle)

	// Verify in database
	dbPage, err := repo.GetByID(ctx, page.ID)
	require.NoError(t, err)
	var dbBlocks []domain.Block
	err = json.Unmarshal(dbPage.Blocks, &dbBlocks)
	require.NoError(t, err)
	assert.Len(t, dbBlocks, 1)
	assert.Equal(t, "hero", dbBlocks[0].Type)
}

func TestPageHandler_CreatePage_WithMeta(t *testing.T) {
	app, db := setupTestApp(t)

	metaData := map[string]interface{}{
		"meta_title":       "Custom SEO Title",
		"meta_description": "Custom SEO description for search engines",
		"og_image":         "/uploads/images/test.jpg",
		"no_index":         true,
	}

	reqBody := CreatePageRequest{
		Slug:   "seo-test-page",
		Title:  "SEO Test Page",
		Status: "draft",
		Meta:   metaData,
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
	assert.Equal(t, "seo-test-page", page.Slug)
	assert.Equal(t, "SEO Test Page", page.Title)

	// Verify meta was saved
	require.NotNil(t, page.Meta)
	var savedMeta map[string]interface{}
	err = json.Unmarshal(page.Meta, &savedMeta)
	require.NoError(t, err)
	assert.Equal(t, "Custom SEO Title", savedMeta["meta_title"])
	assert.Equal(t, "Custom SEO description for search engines", savedMeta["meta_description"])
	assert.Equal(t, "/uploads/images/test.jpg", savedMeta["og_image"])
	assert.Equal(t, true, savedMeta["no_index"])

	// Verify in database
	repo := repository.NewPageRepository(db)
	ctx := database.SetDBInContext(context.Background(), db)
	dbPage, err := repo.GetByID(ctx, page.ID)
	require.NoError(t, err)

	var dbMeta map[string]interface{}
	err = json.Unmarshal(dbPage.Meta, &dbMeta)
	require.NoError(t, err)
	assert.Equal(t, "Custom SEO Title", dbMeta["meta_title"])
	assert.Equal(t, true, dbMeta["no_index"])
}

func TestPageHandler_UpdatePage_WithMeta(t *testing.T) {
	app, db := setupTestApp(t)

	// Create a page first
	repo := repository.NewPageRepository(db)
	ctx := database.SetDBInContext(context.Background(), db)

	page := &domain.Page{
		Slug:   "meta-update-test",
		Title:  "Meta Update Test Page",
		Status: domain.PageStatusDraft,
	}
	err := repo.Create(ctx, page)
	require.NoError(t, err)

	// Update the page with meta
	metaData := map[string]interface{}{
		"meta_title":       "Updated SEO Title",
		"meta_description": "Updated SEO description",
		"og_image":         "/uploads/images/updated.jpg",
		"no_index":         false,
	}

	updateReq := UpdatePageRequest{
		Title: "Updated Title",
		Meta:  metaData,
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

	// Verify meta was updated
	require.NotNil(t, updatedPage.Meta)
	var savedMeta map[string]interface{}
	err = json.Unmarshal(updatedPage.Meta, &savedMeta)
	require.NoError(t, err)
	assert.Equal(t, "Updated SEO Title", savedMeta["meta_title"])
	assert.Equal(t, "Updated SEO description", savedMeta["meta_description"])
	assert.Equal(t, "/uploads/images/updated.jpg", savedMeta["og_image"])
	assert.Equal(t, false, savedMeta["no_index"])

	// Verify in database
	dbPage, err := repo.GetByID(ctx, page.ID)
	require.NoError(t, err)

	var dbMeta map[string]interface{}
	err = json.Unmarshal(dbPage.Meta, &dbMeta)
	require.NoError(t, err)
	assert.Equal(t, "Updated SEO Title", dbMeta["meta_title"])
	assert.Equal(t, false, dbMeta["no_index"])
}

func TestPageHandler_CreatePage_WithEmptyMeta(t *testing.T) {
	app, _ := setupTestApp(t)

	reqBody := CreatePageRequest{
		Slug:   "empty-meta-test",
		Title:  "Empty Meta Test Page",
		Status: "draft",
		Meta:   nil,
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
	assert.Equal(t, "empty-meta-test", page.Slug)

	// Meta should be empty or null
	if len(page.Meta) > 0 {
		var meta map[string]interface{}
		err = json.Unmarshal(page.Meta, &meta)
		require.NoError(t, err)
		assert.Empty(t, meta)
	}
}
