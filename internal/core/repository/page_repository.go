package repository

import (
	"context"

	"gohac/internal/core/domain"

	"github.com/google/uuid"
)

// PageRepository defines the interface for page data access
// This follows the Repository pattern from Clean Architecture
// Implementations will be in internal/adapter/repository
type PageRepository interface {
	// Create creates a new page
	Create(ctx context.Context, page *domain.Page) error

	// GetByID retrieves a page by its UUID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Page, error)

	// GetBySlug retrieves a page by its slug
	GetBySlug(ctx context.Context, slug string) (*domain.Page, error)

	// Update updates an existing page
	Update(ctx context.Context, page *domain.Page) error

	// Delete soft-deletes a page (sets DeletedAt)
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves pages with pagination and filtering
	List(ctx context.Context, opts ListPageOptions) ([]*domain.Page, int64, error)

	// Publish publishes a page (sets status to published and PublishedAt)
	Publish(ctx context.Context, id uuid.UUID) error

	// Unpublish unpublishes a page (sets status to draft)
	Unpublish(ctx context.Context, id uuid.UUID) error
}

// ListPageOptions defines options for listing pages
type ListPageOptions struct {
	Limit    int
	Offset   int
	Status   string // Filter by status (draft, published, archived)
	Search   string // Search in title and description
	TenantID *uuid.UUID
}
