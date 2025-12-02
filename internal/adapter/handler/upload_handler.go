package handler

import (
	"os"

	"gohac/internal/adapter/storage"

	"github.com/gofiber/fiber/v2"
)

// UploadHandler handles file upload operations
type UploadHandler struct {
	storage *storage.Storage
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler() *UploadHandler {
	// Get storage path from environment or use default
	basePath := os.Getenv("STORAGE_PATH")
	if basePath == "" {
		basePath = "./storage"
	}

	// Get base URL from environment or use default
	baseURL := os.Getenv("STORAGE_BASE_URL")
	if baseURL == "" {
		baseURL = "/uploads" // Changed from /static to /uploads for cleaner URLs
	}

	st := storage.NewStorage(basePath, baseURL)

	return &UploadHandler{
		storage: st,
	}
}

// UploadFile handles POST /api/v1/upload
func (h *UploadHandler) UploadFile(c *fiber.Ctx) error {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file provided",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
			"code":  fiber.StatusInternalServerError,
		})
	}
	defer src.Close()

	// Save file
	fileURL, err := h.storage.SaveFile(src, file.Filename, file.Header.Get("Content-Type"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"url": fileURL,
	})
}

// DownloadFromURL handles POST /api/v1/upload/from-url
func (h *UploadHandler) DownloadFromURL(c *fiber.Ctx) error {
	var req struct {
		URL string `json:"url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	if req.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Download and save
	fileURL, err := h.storage.DownloadAndSave(req.URL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusBadRequest,
		})
	}

	return c.JSON(fiber.Map{
		"url": fileURL,
	})
}
