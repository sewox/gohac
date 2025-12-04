package repository

import (
	"context"

	"gohac/internal/core/domain"

	"github.com/google/uuid"
)

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	// Create creates a new category
	Create(ctx context.Context, category *domain.Category) error

	// GetByID retrieves a category by its UUID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)

	// GetBySlug retrieves a category by its slug
	GetBySlug(ctx context.Context, slug string) (*domain.Category, error)

	// Update updates an existing category
	Update(ctx context.Context, category *domain.Category) error

	// Delete deletes a category by its UUID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves a list of categories with pagination
	List(ctx context.Context, limit, offset int) ([]*domain.Category, int64, error)
}
