package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gohac/internal/core/domain"
	"gohac/internal/core/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// pageRepository implements the PageRepository interface using GORM
type pageRepository struct {
	db *gorm.DB
}

// NewPageRepository creates a new page repository instance
func NewPageRepository(db *gorm.DB) repository.PageRepository {
	return &pageRepository{db: db}
}

// Create creates a new page
func (r *pageRepository) Create(ctx context.Context, page *domain.Page) error {
	if page.ID == uuid.Nil {
		page.ID = uuid.New()
	}

	if err := r.db.WithContext(ctx).Create(page).Error; err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}
	return nil
}

// GetByID retrieves a page by its UUID
func (r *pageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Page, error) {
	var page domain.Page
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&page).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("page not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get page: %w", err)
	}
	return &page, nil
}

// GetBySlug retrieves a page by its slug
func (r *pageRepository) GetBySlug(ctx context.Context, slug string) (*domain.Page, error) {
	var page domain.Page
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&page).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("page not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get page by slug: %w", err)
	}
	return &page, nil
}

// Update updates an existing page
func (r *pageRepository) Update(ctx context.Context, page *domain.Page) error {
	if err := r.db.WithContext(ctx).Save(page).Error; err != nil {
		return fmt.Errorf("failed to update page: %w", err)
	}
	return nil
}

// Delete soft-deletes a page
func (r *pageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Page{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete page: %w", err)
	}
	return nil
}

// List retrieves pages with pagination and filtering
func (r *pageRepository) List(ctx context.Context, opts repository.ListPageOptions) ([]*domain.Page, int64, error) {
	var pages []*domain.Page
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Page{})

	// Apply filters
	if opts.Status != "" {
		query = query.Where("status = ?", opts.Status)
	}

	if opts.Search != "" {
		searchTerm := "%" + strings.ToLower(opts.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(slug) LIKE ?", searchTerm, searchTerm)
	}

	if opts.TenantID != nil {
		query = query.Where("tenant_id = ?", opts.TenantID.String())
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count pages: %w", err)
	}

	// Apply pagination
	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		query = query.Offset(opts.Offset)
	}

	// Order by updated_at descending
	query = query.Order("updated_at DESC")

	// Execute query
	if err := query.Find(&pages).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list pages: %w", err)
	}

	return pages, total, nil
}

// Publish publishes a page
func (r *pageRepository) Publish(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&domain.Page{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       domain.PageStatusPublished,
			"published_at": &now,
		}).Error; err != nil {
		return fmt.Errorf("failed to publish page: %w", err)
	}
	return nil
}

// Unpublish unpublishes a page
func (r *pageRepository) Unpublish(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&domain.Page{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       domain.PageStatusDraft,
			"published_at": nil,
		}).Error; err != nil {
		return fmt.Errorf("failed to unpublish page: %w", err)
	}
	return nil
}
