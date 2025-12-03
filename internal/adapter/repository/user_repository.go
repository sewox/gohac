package repository

import (
	"context"
	"fmt"

	"gohac/internal/core/domain"
	"gohac/internal/core/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// userRepository implements the UserRepository interface using GORM
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by its UUID (accepts string or uuid.UUID)
func (r *userRepository) GetByID(ctx context.Context, id interface{}) (*domain.User, error) {
	var user domain.User
	var err error
	var userUUID uuid.UUID

	// Handle both string and uuid.UUID types
	switch v := id.(type) {
	case string:
		parsedUUID, parseErr := uuid.Parse(v)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid UUID format: %w", parseErr)
		}
		userUUID = parsedUUID
	case uuid.UUID:
		userUUID = v
	default:
		return nil, fmt.Errorf("invalid ID type, expected string or uuid.UUID")
	}

	err = r.db.WithContext(ctx).First(&user, "id = ?", userUUID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete deletes a user by its UUID (accepts string or uuid.UUID)
func (r *userRepository) Delete(ctx context.Context, id interface{}) error {
	var userUUID uuid.UUID

	// Handle both string and uuid.UUID types
	switch v := id.(type) {
	case string:
		parsedUUID, parseErr := uuid.Parse(v)
		if parseErr != nil {
			return fmt.Errorf("invalid UUID format: %w", parseErr)
		}
		userUUID = parsedUUID
	case uuid.UUID:
		userUUID = v
	default:
		return fmt.Errorf("invalid ID type, expected string or uuid.UUID")
	}

	if err := r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", userUUID).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List retrieves a list of users with pagination
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Fetch paginated records
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}
