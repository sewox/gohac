package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostStatus defines the status of a post
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

// Post represents a blog post
type Post struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TenantID      string     `gorm:"index" json:"tenant_id"` // Empty string for community edition
	Title         string     `gorm:"type:varchar(255);not null" json:"title"`
	Slug          string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"slug"`
	Excerpt       string     `gorm:"type:text" json:"excerpt"`
	Content       string     `gorm:"type:text" json:"content"` // JSON Blocks array
	FeaturedImage string     `gorm:"type:varchar(500)" json:"featured_image"`
	Status        PostStatus `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	PublishedAt   *time.Time `json:"published_at"`
	AuthorID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"author_id"`
	Author        User       `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Categories    []Category `gorm:"many2many:post_categories;" json:"categories,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// BeforeCreate is a GORM hook that generates UUID before creating a post
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for GORM
func (Post) TableName() string {
	return "posts"
}

// Category represents a blog category/taxonomy
type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	TenantID    string    `gorm:"index" json:"tenant_id"` // Empty string for community edition
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Posts       []Post    `gorm:"many2many:post_categories;" json:"posts,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate is a GORM hook that generates UUID before creating a category
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for GORM
func (Category) TableName() string {
	return "categories"
}
