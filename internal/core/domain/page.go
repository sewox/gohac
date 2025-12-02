package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Page represents a content page in the CMS
// Uses JSONB blocks for flexible, schema-less content structure
type Page struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	TenantID    string         `gorm:"index;not null" json:"tenant_id"` // Empty string for community edition
	Slug        string         `gorm:"index;not null" json:"slug"`
	Title       string         `gorm:"not null" json:"title"`
	Blocks      datatypes.JSON `gorm:"type:jsonb" json:"blocks"` // Array of Block objects
	Status      PageStatus     `gorm:"type:varchar(20);default:'draft'" json:"status"`
	Meta        datatypes.JSON `gorm:"type:jsonb" json:"meta"` // SEO, custom fields, etc.
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	PublishedAt *time.Time     `json:"published_at,omitempty"`
}

// BeforeCreate is a GORM hook that generates UUID before creating a page
// This ensures SQLite compatibility (SQLite doesn't have gen_random_uuid())
func (p *Page) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// PageStatus represents the publication status of a page
type PageStatus string

const (
	PageStatusDraft     PageStatus = "draft"
	PageStatusPublished PageStatus = "published"
	PageStatusArchived  PageStatus = "archived"
)

// TableName specifies the table name for GORM
func (Page) TableName() string {
	return "pages"
}
