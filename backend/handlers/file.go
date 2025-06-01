package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

// UploadFile handles file uploads
func (h *Handler) UploadFile(c echo.Context) error {
	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No file uploaded",
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open file",
		})
	}
	defer src.Close()

	// Create uploads directory if it doesn't exist
	uploadsDir := "../uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create uploads directory",
		})
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	filepath := filepath.Join(uploadsDir, filename)

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create file",
		})
	}
	defer dst.Close()

	// Copy file
	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save file",
		})
	}

	// Return file information
	return c.JSON(http.StatusOK, map[string]interface{}{
		"filename": filename,
		"size":     file.Size,
		"url":      fmt.Sprintf("/api/files/%s", filename),
	})
}

// GetFile serves uploaded files
func (h *Handler) GetFile(c echo.Context) error {
	filename := c.Param("id")
	filepath := filepath.Join("../uploads", filename)

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "File not found",
		})
	}

	return c.File(filepath)
}