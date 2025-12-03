package handler

import (
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// MediaHandler handles media-related HTTP requests
type MediaHandler struct {
	uploadPath string
}

// NewMediaHandler creates a new media handler instance
func NewMediaHandler() *MediaHandler {
	uploadPath := os.Getenv("STORAGE_PATH")
	if uploadPath == "" {
		uploadPath = "./storage/uploads"
	}
	return &MediaHandler{
		uploadPath: uploadPath,
	}
}

// MediaItem represents a media file
type MediaItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}

// ListMedia handles GET /api/v1/media (protected endpoint)
func (h *MediaHandler) ListMedia(c *fiber.Ctx) error {
	var mediaItems []MediaItem

	// Read directory
	entries, err := os.ReadDir(h.uploadPath)
	if err != nil {
		// If directory doesn't exist, return empty list
		if os.IsNotExist(err) {
			return c.JSON(fiber.Map{
				"data":  []MediaItem{},
				"total": 0,
			})
		}
		log.Printf("Error reading media directory: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list media files",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Process each file
	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip directories
		}

		fileInfo, err := entry.Info()
		if err != nil {
			log.Printf("Error getting file info for %s: %v", entry.Name(), err)
			continue
		}

		// Determine MIME type
		ext := filepath.Ext(entry.Name())
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		// Build URL (relative to the static file serving path)
		url := "/uploads/" + entry.Name()

		mediaItems = append(mediaItems, MediaItem{
			Name: entry.Name(),
			URL:  url,
			Size: fileInfo.Size(),
			Type: mimeType,
		})
	}

	return c.JSON(fiber.Map{
		"data":  mediaItems,
		"total": len(mediaItems),
	})
}

// GetMediaInfo handles GET /api/v1/media/:filename (protected endpoint)
func (h *MediaHandler) GetMediaInfo(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Filename is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Sanitize filename to prevent directory traversal
	filename = filepath.Base(filename)
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid filename",
			"code":  fiber.StatusBadRequest,
		})
	}

	filePath := filepath.Join(h.uploadPath, filename)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "File not found",
				"code":  fiber.StatusNotFound,
			})
		}
		log.Printf("Error getting file info: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get file info",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Determine MIME type
	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	url := "/uploads/" + filename

	return c.JSON(MediaItem{
		Name: filename,
		URL:  url,
		Size: fileInfo.Size(),
		Type: mimeType,
	})
}
