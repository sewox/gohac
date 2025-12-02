package repository

import (
	"context"
	"fmt"

	"gohac/internal/core/domain"
	"gohac/internal/core/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// menuRepository implements the MenuRepository interface using GORM
type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository creates a new menu repository instance
func NewMenuRepository(db *gorm.DB) repository.MenuRepository {
	return &menuRepository{db: db}
}

// Create creates a new menu
func (r *menuRepository) Create(ctx context.Context, menu *domain.Menu) error {
	if err := r.db.WithContext(ctx).Create(menu).Error; err != nil {
		return fmt.Errorf("failed to create menu: %w", err)
	}
	return nil
}

// GetByID retrieves a menu by its UUID
func (r *menuRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Menu, error) {
	var menu domain.Menu
	err := r.db.WithContext(ctx).First(&menu, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("menu not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get menu: %w", err)
	}
	return &menu, nil
}

// Update updates an existing menu
func (r *menuRepository) Update(ctx context.Context, menu *domain.Menu) error {
	if err := r.db.WithContext(ctx).Save(menu).Error; err != nil {
		return fmt.Errorf("failed to update menu: %w", err)
	}
	return nil
}

// Delete deletes a menu by ID
func (r *menuRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Menu{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete menu: %w", err)
	}
	return nil
}

// List retrieves all menus with pagination
func (r *menuRepository) List(ctx context.Context, limit, offset int) ([]*domain.Menu, int64, error) {
	var menus []*domain.Menu
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.Menu{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count menus: %w", err)
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&menus).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list menus: %w", err)
	}

	return menus, total, nil
}
