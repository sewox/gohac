package repository

import (
	"context"
	"fmt"

	"gohac/internal/core/domain"
	"gohac/internal/core/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// postRepository implements the PostRepository interface using GORM
type postRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new post repository instance
func NewPostRepository(db *gorm.DB) repository.PostRepository {
	return &postRepository{db: db}
}

// Create creates a new post
func (r *postRepository) Create(ctx context.Context, post *domain.Post) error {
	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

// GetByID retrieves a post by its UUID
func (r *postRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	var post domain.Post
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Categories").
		First(&post, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("post not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return &post, nil
}

// GetBySlug retrieves a post by its slug
func (r *postRepository) GetBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	var post domain.Post
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Categories").
		Where("slug = ? AND status = ?", slug, domain.PostStatusPublished).
		First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("post not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get post by slug: %w", err)
	}
	return &post, nil
}

// Update updates an existing post
func (r *postRepository) Update(ctx context.Context, post *domain.Post) error {
	if err := r.db.WithContext(ctx).Save(post).Error; err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

// Delete deletes a post by its UUID
func (r *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Post{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

// List retrieves a list of posts with pagination
func (r *postRepository) List(ctx context.Context, limit, offset int, status *domain.PostStatus) ([]*domain.Post, int64, error) {
	var posts []*domain.Post
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Post{})

	// Filter by status if provided
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Fetch paginated records
	query = query.Preload("Author").Preload("Categories").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC")

	if err := query.Find(&posts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list posts: %w", err)
	}

	return posts, total, nil
}

// ListByCategory retrieves posts by category ID
func (r *postRepository) ListByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*domain.Post, int64, error) {
	var posts []*domain.Post
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).
		Model(&domain.Post{}).
		Joins("JOIN post_categories ON posts.id = post_categories.post_id").
		Where("post_categories.category_id = ?", categoryID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count posts by category: %w", err)
	}

	// Fetch paginated records
	if err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Categories").
		Joins("JOIN post_categories ON posts.id = post_categories.post_id").
		Where("post_categories.category_id = ?", categoryID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list posts by category: %w", err)
	}

	return posts, total, nil
}
