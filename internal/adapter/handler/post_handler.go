package handler

import (
	"log"
	"strconv"
	"strings"
	"time"

	"gohac/internal/adapter/database"
	"gohac/internal/adapter/repository"
	"gohac/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostHandler handles post-related HTTP requests
type PostHandler struct {
	db *gorm.DB
}

// NewPostHandler creates a new post handler instance
func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{
		db: db,
	}
}

// CreatePostRequest represents the request body for creating a post
type CreatePostRequest struct {
	Title         string   `json:"title" validate:"required"`
	Slug          string   `json:"slug" validate:"required"`
	Excerpt       string   `json:"excerpt"`
	Content       string   `json:"content"` // JSON Blocks array
	FeaturedImage string   `json:"featured_image"`
	Status        string   `json:"status" validate:"required,oneof=draft published archived"`
	CategoryIDs   []string `json:"category_ids"`
}

// UpdatePostRequest represents the request body for updating a post
type UpdatePostRequest struct {
	Title         string   `json:"title,omitempty"`
	Slug          string   `json:"slug,omitempty"`
	Excerpt       string   `json:"excerpt,omitempty"`
	Content       string   `json:"content,omitempty"`
	FeaturedImage string   `json:"featured_image,omitempty"`
	Status        string   `json:"status,omitempty"`
	CategoryIDs   []string `json:"category_ids,omitempty"`
}

// CreatePost handles POST /api/v1/posts (protected endpoint)
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	var req CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Validate required fields
	if req.Title == "" || req.Slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title and slug are required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get current user ID
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
			"code":  fiber.StatusUnauthorized,
		})
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user ID format",
			"code":  fiber.StatusInternalServerError,
		})
	}

	authorUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user ID format",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Get database from context
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	// Validate status
	status := domain.PostStatus(strings.ToLower(req.Status))
	if status != domain.PostStatusDraft && status != domain.PostStatusPublished && status != domain.PostStatusArchived {
		status = domain.PostStatusDraft
	}

	// Create post
	post := &domain.Post{
		Title:         req.Title,
		Slug:          req.Slug,
		Excerpt:       req.Excerpt,
		Content:       req.Content,
		FeaturedImage: req.FeaturedImage,
		Status:        status,
		AuthorID:      authorUUID,
	}

	// Set published_at if status is published
	if status == domain.PostStatusPublished {
		now := time.Now()
		post.PublishedAt = &now
	}

	// Load categories if provided
	if len(req.CategoryIDs) > 0 {
		categoryRepo := repository.NewCategoryRepository(db)
		var categories []domain.Category
		for _, catIDStr := range req.CategoryIDs {
			catID, err := uuid.Parse(catIDStr)
			if err != nil {
				continue
			}
			category, err := categoryRepo.GetByID(c.Context(), catID)
			if err == nil {
				categories = append(categories, *category)
			}
		}
		post.Categories = categories
	}

	postRepo := repository.NewPostRepository(db)
	if err := postRepo.Create(c.Context(), post); err != nil {
		log.Printf("Error creating post: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create post",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Reload post with relations
	post, err = postRepo.GetByID(c.Context(), post.ID)
	if err != nil {
		log.Printf("Error reloading post: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}

// ListPosts handles GET /api/v1/posts (protected endpoint)
func (h *PostHandler) ListPosts(c *fiber.Ctx) error {
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	postRepo := repository.NewPostRepository(db)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	statusStr := c.Query("status")

	var status *domain.PostStatus
	if statusStr != "" {
		s := domain.PostStatus(strings.ToLower(statusStr))
		status = &s
	}

	posts, total, err := postRepo.List(c.Context(), limit, offset, status)
	if err != nil {
		log.Printf("Error listing posts: %v", err)
		// If table doesn't exist, return empty list instead of error
		if strings.Contains(err.Error(), "no such table") || strings.Contains(err.Error(), "doesn't exist") {
			return c.JSON(fiber.Map{
				"data":   []*domain.Post{},
				"total":  0,
				"limit":  limit,
				"offset": offset,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list posts",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"data":   posts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetPost handles GET /api/v1/posts/:id (protected endpoint)
func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID format",
			"code":  fiber.StatusBadRequest,
		})
	}

	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	postRepo := repository.NewPostRepository(db)
	post, err := postRepo.GetByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "post not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Post not found",
				"code":  fiber.StatusNotFound,
			})
		}
		log.Printf("Error getting post: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get post",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(post)
}

// UpdatePost handles PUT /api/v1/posts/:id (protected endpoint)
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID format",
			"code":  fiber.StatusBadRequest,
		})
	}

	var req UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	postRepo := repository.NewPostRepository(db)
	post, err := postRepo.GetByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "post not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Post not found",
				"code":  fiber.StatusNotFound,
			})
		}
		log.Printf("Error getting post for update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get post",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Update fields
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Slug != "" {
		post.Slug = req.Slug
	}
	if req.Excerpt != "" || req.Excerpt == "" {
		post.Excerpt = req.Excerpt
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.FeaturedImage != "" {
		post.FeaturedImage = req.FeaturedImage
	}
	if req.Status != "" {
		status := domain.PostStatus(strings.ToLower(req.Status))
		if status == domain.PostStatusDraft || status == domain.PostStatusPublished || status == domain.PostStatusArchived {
			oldStatus := post.Status
			post.Status = status
			// Set published_at if transitioning to published
			if status == domain.PostStatusPublished && oldStatus != domain.PostStatusPublished {
				now := time.Now()
				post.PublishedAt = &now
			}
		}
	}

	// Update categories if provided
	if req.CategoryIDs != nil {
		categoryRepo := repository.NewCategoryRepository(db)
		var categories []domain.Category
		for _, catIDStr := range req.CategoryIDs {
			catID, err := uuid.Parse(catIDStr)
			if err != nil {
				continue
			}
			category, err := categoryRepo.GetByID(c.Context(), catID)
			if err == nil {
				categories = append(categories, *category)
			}
		}
		post.Categories = categories
	}

	if err := postRepo.Update(c.Context(), post); err != nil {
		log.Printf("Error updating post: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update post",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Reload post with relations
	post, err = postRepo.GetByID(c.Context(), post.ID)
	if err != nil {
		log.Printf("Error reloading post: %v", err)
	}

	return c.JSON(post)
}

// DeletePost handles DELETE /api/v1/posts/:id (protected endpoint)
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID format",
			"code":  fiber.StatusBadRequest,
		})
	}

	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	postRepo := repository.NewPostRepository(db)
	if err := postRepo.Delete(c.Context(), id); err != nil {
		log.Printf("Error deleting post: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete post",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// GetPostBySlugPublic handles GET /api/public/posts/:slug (public endpoint)
func (h *PostHandler) GetPostBySlugPublic(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Slug is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	postRepo := repository.NewPostRepository(db)
	post, err := postRepo.GetBySlug(c.Context(), slug)
	if err != nil {
		if strings.Contains(err.Error(), "post not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Post not found",
				"code":  fiber.StatusNotFound,
			})
		}
		log.Printf("Error getting post by slug: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get post",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(post)
}

// ListPostsPublic handles GET /api/public/posts (public endpoint)
func (h *PostHandler) ListPostsPublic(c *fiber.Ctx) error {
	db, err := database.GetDBFromContext(c.Context())
	if err != nil {
		db = h.db
	}

	postRepo := repository.NewPostRepository(db)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	// Only show published posts
	status := domain.PostStatusPublished
	posts, total, err := postRepo.List(c.Context(), limit, offset, &status)
	if err != nil {
		log.Printf("Error listing posts: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list posts",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"data":   posts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}
