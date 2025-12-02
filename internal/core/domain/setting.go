package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// GlobalSettings represents site-wide configuration
type GlobalSettings struct {
	SiteName     string `json:"site_name"`
	Logo         string `json:"logo"`
	Favicon      string `json:"favicon"`
	ContactEmail string `json:"contact_email"`
}

// MenuItem represents a single menu item (can be nested)
type MenuItem struct {
	Label    string     `json:"label"`
	URL      string     `json:"url"`
	Target   string     `json:"target,omitempty"` // "_blank", "_self", etc.
	Children []MenuItem `json:"children,omitempty"`
}

// Menu represents a navigation menu (reusable, can be used anywhere)
type Menu struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	TenantID    string         `gorm:"index" json:"tenant_id"`                 // Empty string for community edition
	Name        string         `gorm:"type:varchar(100);not null" json:"name"` // Menu name (e.g., "Main Navigation", "Footer Links")
	Description string         `gorm:"type:text" json:"description,omitempty"` // Optional description
	Items       datatypes.JSON `gorm:"type:jsonb" json:"items"`                // Array of MenuItem objects
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// BeforeCreate is a GORM hook that generates UUID before creating a menu
func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for GORM
func (Menu) TableName() string {
	return "menus"
}

// SystemConfig represents a key-value configuration stored in the database
// Used for storing GlobalSettings as a single JSON blob
type SystemConfig struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	TenantID  string         `gorm:"index;not null" json:"tenant_id"` // Empty string for community edition
	Key       string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_tenant_key" json:"key"`
	Value     datatypes.JSON `gorm:"type:jsonb" json:"value"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// BeforeCreate is a GORM hook that generates UUID before creating a system config
func (s *SystemConfig) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for GORM
func (SystemConfig) TableName() string {
	return "system_configs"
}
