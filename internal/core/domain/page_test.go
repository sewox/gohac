package domain

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPage_BeforeCreate_GeneratesUUID(t *testing.T) {
	// Test Case 1: Create a Page with an empty ID. Assert that BeforeCreate generates a valid UUID.
	// Use a real in-memory DB to properly test the hook
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate schema
	err = db.AutoMigrate(&Page{})
	require.NoError(t, err)

	page := &Page{
		ID:       uuid.Nil, // Empty UUID
		TenantID: "",
		Slug:     "test-page",
		Title:    "Test Page",
		Status:   PageStatusDraft,
	}

	// Create will trigger BeforeCreate hook
	err = db.Create(page).Error
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, page.ID)
	assert.NotEqual(t, uuid.UUID{}, page.ID)
}

func TestPage_BeforeCreate_PreservesExistingID(t *testing.T) {
	// Test Case 2: Create a Page with a specific ID. Assert that the ID is NOT overwritten.
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate schema
	err = db.AutoMigrate(&Page{})
	require.NoError(t, err)

	existingID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	page := &Page{
		ID:       existingID,
		TenantID: "",
		Slug:     "test-page",
		Title:    "Test Page",
		Status:   PageStatusDraft,
	}

	// Create will trigger BeforeCreate hook
	err = db.Create(page).Error
	require.NoError(t, err)
	assert.Equal(t, existingID, page.ID)
}

func TestPage_Blocks_JSONMarshaling(t *testing.T) {
	// Test Case 3: Verify Blocks JSON marshaling/unmarshaling works correctly with the Block struct.
	blocks := []Block{
		{
			ID:   uuid.New().String(),
			Type: string(BlockTypeHero),
			Data: json.RawMessage(`{"title":"Hero Title","subtitle":"Hero Subtitle"}`),
		},
		{
			ID:   uuid.New().String(),
			Type: string(BlockTypeText),
			Data: json.RawMessage(`{"content":"<p>Hello World</p>","align":"left"}`),
		},
	}

	// Marshal blocks to JSON
	blocksJSON, err := json.Marshal(blocks)
	require.NoError(t, err)

	// Create a page with blocks
	page := &Page{
		ID:       uuid.New(),
		TenantID: "",
		Slug:     "test-page",
		Title:    "Test Page",
		Blocks:   datatypes.JSON(blocksJSON),
		Status:   PageStatusDraft,
	}

	// Verify blocks can be unmarshaled
	var unmarshaledBlocks []Block
	err = json.Unmarshal(page.Blocks, &unmarshaledBlocks)
	require.NoError(t, err)

	assert.Len(t, unmarshaledBlocks, 2)
	assert.Equal(t, blocks[0].ID, unmarshaledBlocks[0].ID)
	assert.Equal(t, blocks[0].Type, unmarshaledBlocks[0].Type)
	assert.Equal(t, blocks[1].ID, unmarshaledBlocks[1].ID)
	assert.Equal(t, blocks[1].Type, unmarshaledBlocks[1].Type)

	// Verify block data can be parsed
	var heroData HeroBlockData
	err = json.Unmarshal(unmarshaledBlocks[0].Data, &heroData)
	require.NoError(t, err)
	assert.Equal(t, "Hero Title", heroData.Title)
	assert.Equal(t, "Hero Subtitle", heroData.Subtitle)

	var textData TextBlockData
	err = json.Unmarshal(unmarshaledBlocks[1].Data, &textData)
	require.NoError(t, err)
	assert.Equal(t, "<p>Hello World</p>", textData.Content)
	assert.Equal(t, "left", textData.Align)
}
