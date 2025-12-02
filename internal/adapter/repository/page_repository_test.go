package repository

import (
	"context"
	"testing"

	"gohac/internal/core/domain"
	"gohac/internal/core/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
// Each test gets its own isolated database
func setupTestDB(t *testing.T) *gorm.DB {
	// Use unique memory database for each test to avoid interference
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate schema
	err = db.AutoMigrate(&domain.Page{})
	require.NoError(t, err)

	return db
}

func TestPageRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPageRepository(db)
	ctx := context.Background()

	page := &domain.Page{
		TenantID: "",
		Slug:     "test-page",
		Title:    "Test Page",
		Status:   domain.PageStatusDraft,
	}

	err := repo.Create(ctx, page)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, page.ID)
}

func TestPageRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPageRepository(db)
	ctx := context.Background()

	// Create a page first
	page := &domain.Page{
		TenantID: "",
		Slug:     "test-page",
		Title:    "Test Page",
		Status:   domain.PageStatusDraft,
	}
	err := repo.Create(ctx, page)
	require.NoError(t, err)

	// Retrieve by ID
	retrieved, err := repo.GetByID(ctx, page.ID)
	require.NoError(t, err)
	assert.Equal(t, page.ID, retrieved.ID)
	assert.Equal(t, page.Slug, retrieved.Slug)
	assert.Equal(t, page.Title, retrieved.Title)
}

func TestPageRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPageRepository(db)
	ctx := context.Background()

	nonExistentID := uuid.New()
	_, err := repo.GetByID(ctx, nonExistentID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "page not found")
}

func TestPageRepository_GetBySlug(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPageRepository(db)
	ctx := context.Background()

	// Create a page
	page := &domain.Page{
		TenantID: "",
		Slug:     "test-slug",
		Title:    "Test Page",
		Status:   domain.PageStatusDraft,
	}
	err := repo.Create(ctx, page)
	require.NoError(t, err)

	// Retrieve by slug
	retrieved, err := repo.GetBySlug(ctx, "test-slug")
	require.NoError(t, err)
	assert.Equal(t, page.ID, retrieved.ID)
	assert.Equal(t, page.Slug, retrieved.Slug)
}

func TestPageRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPageRepository(db)
	ctx := context.Background()

	// Create a page
	page := &domain.Page{
		TenantID: "",
		Slug:     "test-page",
		Title:    "Original Title",
		Status:   domain.PageStatusDraft,
	}
	err := repo.Create(ctx, page)
	require.NoError(t, err)

	// Update the title
	page.Title = "Updated Title"
	err = repo.Update(ctx, page)
	require.NoError(t, err)

	// Verify update
	updated, err := repo.GetByID(ctx, page.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)
}

func TestPageRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPageRepository(db)
	ctx := context.Background()

	// Create a page
	page := &domain.Page{
		TenantID: "",
		Slug:     "test-page",
		Title:    "Test Page",
		Status:   domain.PageStatusDraft,
	}
	err := repo.Create(ctx, page)
	require.NoError(t, err)

	// Delete the page
	err = repo.Delete(ctx, page.ID)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = repo.GetByID(ctx, page.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "page not found")
}

func TestPageRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPageRepository(db)
	ctx := context.Background()

	// Create multiple pages
	for i := 0; i < 5; i++ {
		page := &domain.Page{
			TenantID: "",
			Slug:     "test-page-" + string(rune('a'+i)),
			Title:    "Test Page " + string(rune('A'+i)),
			Status:   domain.PageStatusDraft,
		}
		err := repo.Create(ctx, page)
		require.NoError(t, err)
	}

	// List all pages
	opts := repository.ListPageOptions{
		Limit:  10,
		Offset: 0,
	}
	pages, total, err := repo.List(ctx, opts)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, pages, 5)
}

func TestPageRepository_List_WithStatusFilter(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPageRepository(db)
	ctx := context.Background()

	// Create pages with different statuses
	draftPage := &domain.Page{
		TenantID: "",
		Slug:     "draft-page",
		Title:    "Draft Page",
		Status:   domain.PageStatusDraft,
	}
	err := repo.Create(ctx, draftPage)
	require.NoError(t, err)

	publishedPage := &domain.Page{
		TenantID: "",
		Slug:     "published-page",
		Title:    "Published Page",
		Status:   domain.PageStatusPublished,
	}
	err = repo.Create(ctx, publishedPage)
	require.NoError(t, err)

	// List only published pages
	opts := repository.ListPageOptions{
		Limit:  10,
		Offset: 0,
		Status: string(domain.PageStatusPublished),
	}
	pages, total, err := repo.List(ctx, opts)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, pages, 1)
	assert.Equal(t, publishedPage.ID, pages[0].ID)
}

func TestPageRepository_Publish(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPageRepository(db)
	ctx := context.Background()

	// Create a draft page
	page := &domain.Page{
		TenantID: "",
		Slug:     "test-page",
		Title:    "Test Page",
		Status:   domain.PageStatusDraft,
	}
	err := repo.Create(ctx, page)
	require.NoError(t, err)

	// Publish the page
	err = repo.Publish(ctx, page.ID)
	require.NoError(t, err)

	// Verify status changed
	updated, err := repo.GetByID(ctx, page.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.PageStatusPublished, updated.Status)
	assert.NotNil(t, updated.PublishedAt)
}

func TestPageRepository_Unpublish(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPageRepository(db)
	ctx := context.Background()

	// Create and publish a page
	page := &domain.Page{
		TenantID: "",
		Slug:     "test-page",
		Title:    "Test Page",
		Status:   domain.PageStatusPublished,
	}
	err := repo.Create(ctx, page)
	require.NoError(t, err)

	// Unpublish the page
	err = repo.Unpublish(ctx, page.ID)
	require.NoError(t, err)

	// Verify status changed
	updated, err := repo.GetByID(ctx, page.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.PageStatusDraft, updated.Status)
	assert.Nil(t, updated.PublishedAt)
}
