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
	BlockTypeHero        BlockType = "hero"
	BlockTypeText        BlockType = "text"
	BlockTypeImage       BlockType = "image"
	BlockTypeGallery     BlockType = "gallery"
	BlockTypeVideo       BlockType = "video"
	BlockTypeQuote       BlockType = "quote"
	BlockTypeCode        BlockType = "code"
	BlockTypeFeatures    BlockType = "features"
	BlockTypePricing     BlockType = "pricing"
	BlockTypeFAQ         BlockType = "faq"
	BlockTypeTestimonial BlockType = "testimonial"
	BlockTypeCTA         BlockType = "cta"
	BlockTypeMenu        BlockType = "menu"
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

// FeaturesBlockData represents data for a features block
type FeaturesBlockData struct {
	Title    string        `json:"title,omitempty"`
	Subtitle string        `json:"subtitle,omitempty"`
	Items    []FeatureItem `json:"items"`
	Columns  int           `json:"columns,omitempty"` // 2, 3, or 4
}

// FeatureItem represents a single feature item
type FeatureItem struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"` // Icon name or URL
}

// PricingBlockData represents data for a pricing block
type PricingBlockData struct {
	Title    string        `json:"title,omitempty"`
	Subtitle string        `json:"subtitle,omitempty"`
	Plans    []PricingPlan `json:"plans"`
}

// PricingPlan represents a single pricing plan
type PricingPlan struct {
	Name        string   `json:"name"`
	Price       string   `json:"price"` // e.g., "$99/month" or "$999/year"
	Description string   `json:"description,omitempty"`
	Features    []string `json:"features"` // List of feature strings
	ButtonText  string   `json:"button_text,omitempty"`
	ButtonURL   string   `json:"button_url,omitempty"`
	Highlighted bool     `json:"highlighted,omitempty"` // For featured plan
}

// FAQBlockData represents data for an FAQ block
type FAQBlockData struct {
	Title string    `json:"title,omitempty"`
	Items []FAQItem `json:"items"`
}

// FAQItem represents a single FAQ item
type FAQItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// TestimonialBlockData represents data for a testimonials block
type TestimonialBlockData struct {
	Title        string            `json:"title,omitempty"`
	Subtitle     string            `json:"subtitle,omitempty"`
	Testimonials []TestimonialItem `json:"testimonials"`
}

// TestimonialItem represents a single testimonial
type TestimonialItem struct {
	Quote     string `json:"quote"`
	Author    string `json:"author"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Role      string `json:"role,omitempty"` // e.g., "CEO, Company Name"
}

// VideoBlockData represents data for a video block
type VideoBlockData struct {
	URL         string `json:"url"` // YouTube, Vimeo, or direct video URL
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Autoplay    bool   `json:"autoplay,omitempty"`
	Loop        bool   `json:"loop,omitempty"`
}

// CTABlockData represents data for a CTA block
type CTABlockData struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle,omitempty"`
	ButtonText  string `json:"button_text"`
	ButtonURL   string `json:"button_url"`
	ButtonStyle string `json:"button_style,omitempty"` // primary, secondary, outline
	Background  string `json:"background,omitempty"`   // Color or gradient
}

// MenuBlockData represents data for a menu block
// MenuID references a Menu entity by UUID
type MenuBlockData struct {
	MenuID string `json:"menu_id"`         // UUID of the menu to display
	Style  string `json:"style,omitempty"` // horizontal, vertical, dropdown
}
