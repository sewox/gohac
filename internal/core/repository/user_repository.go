package repository

import (
	"context"

	"gohac/internal/core/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// GetByID retrieves a user by its UUID (accepts string or uuid.UUID)
	GetByID(ctx context.Context, id interface{}) (*domain.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *domain.User) error

	// Delete deletes a user by its UUID (accepts string or uuid.UUID)
	Delete(ctx context.Context, id interface{}) error

	// List retrieves a list of users with pagination
	List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
}
