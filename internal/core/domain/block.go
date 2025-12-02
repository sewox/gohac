package domain

import "encoding/json"

// Block represents a single content block within a page
// This follows the Block Protocol pattern for flexible content composition
type Block struct {
	ID   string          `json:"id"`   // UUID string
	Type string          `json:"type"` // e.g., "hero", "text", "image", "gallery"
	Data json.RawMessage `json:"data"` // Flexible JSON data specific to block type
}

// BlockType defines common block types
type BlockType string

const (
	BlockTypeHero    BlockType = "hero"
	BlockTypeText    BlockType = "text"
	BlockTypeImage   BlockType = "image"
	BlockTypeGallery BlockType = "gallery"
	BlockTypeVideo   BlockType = "video"
	BlockTypeQuote   BlockType = "quote"
	BlockTypeCode    BlockType = "code"
)

// BlockData represents the structure for common block data types
// This is a helper interface - actual implementations will vary by block type
type BlockData interface {
	Validate() error
}

// HeroBlockData represents data for a hero block
type HeroBlockData struct {
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle,omitempty"`
	ImageURL   string `json:"image_url,omitempty"`
	CTA        *CTA   `json:"cta,omitempty"`
	Background string `json:"background,omitempty"` // Color or gradient
}

// TextBlockData represents data for a text block
type TextBlockData struct {
	Content string `json:"content"`         // HTML or Markdown
	Align   string `json:"align,omitempty"` // left, center, right
}

// ImageBlockData represents data for an image block
type ImageBlockData struct {
	URL     string `json:"url"`
	Alt     string `json:"alt"`
	Caption string `json:"caption,omitempty"`
	Width   int    `json:"width,omitempty"`
	Height  int    `json:"height,omitempty"`
}

// CTA (Call To Action) represents a button or link
type CTA struct {
	Text string `json:"text"`
	URL  string `json:"url"`
	Type string `json:"type,omitempty"` // primary, secondary, etc.
}
