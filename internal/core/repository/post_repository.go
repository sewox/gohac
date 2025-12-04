package repository

import (
	"context"

	"gohac/internal/core/domain"

	"github.com/google/uuid"
)

// PostRepository defines the interface for post data access
type PostRepository interface {
	// Create creates a new post
	Create(ctx context.Context, post *domain.Post) error

	// GetByID retrieves a post by its UUID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Post, error)

	// GetBySlug retrieves a post by its slug
	GetBySlug(ctx context.Context, slug string) (*domain.Post, error)

	// Update updates an existing post
	Update(ctx context.Context, post *domain.Post) error

	// Delete deletes a post by its UUID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves a list of posts with pagination
	List(ctx context.Context, limit, offset int, status *domain.PostStatus) ([]*domain.Post, int64, error)

	// ListByCategory retrieves posts by category ID
	ListByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*domain.Post, int64, error)
}
