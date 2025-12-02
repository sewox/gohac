package repository

import (
	"context"

	"gohac/internal/core/domain"

	"github.com/google/uuid"
)

// MenuRepository defines the interface for menu data access
type MenuRepository interface {
	// Create creates a new menu
	Create(ctx context.Context, menu *domain.Menu) error

	// GetByID retrieves a menu by its UUID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Menu, error)

	// Update updates an existing menu
	Update(ctx context.Context, menu *domain.Menu) error

	// Delete deletes a menu by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all menus with pagination
	List(ctx context.Context, limit, offset int) ([]*domain.Menu, int64, error)
}
