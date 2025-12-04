package repository

import (
	"context"
	"fmt"

	"gohac/internal/core/domain"
	"gohac/internal/core/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// categoryRepository implements the CategoryRepository interface using GORM
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository instance
func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepository{db: db}
}

// Create creates a new category
func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

// GetByID retrieves a category by its UUID
func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	var category domain.Category
	err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

// GetBySlug retrieves a category by its slug
func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}
	return &category, nil
}

// Update updates an existing category
func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

// Delete deletes a category by its UUID
func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Category{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

// List retrieves a list of categories with pagination
func (r *categoryRepository) List(ctx context.Context, limit, offset int) ([]*domain.Category, int64, error) {
	var categories []*domain.Category
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&domain.Category{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	// Fetch paginated records
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&categories).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, total, nil
}
