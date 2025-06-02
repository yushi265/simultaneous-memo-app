package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
)

// ServeImage serves images with optional resizing
func (h *Handler) ServeImage(c echo.Context) error {
	// Get image path
	path := c.Param("*")
	imagePath := filepath.Join("../uploads", path)

	// Security check
	if strings.Contains(path, "..") {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "不正なパスです",
		})
	}

	// Parse query parameters for resizing
	widthStr := c.QueryParam("w")
	heightStr := c.QueryParam("h")
	quality := c.QueryParam("q")
	format := c.QueryParam("format")

	// Check if original file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "画像が見つかりません",
		})
	}

	// If no resizing parameters, serve original
	if widthStr == "" && heightStr == "" && quality == "" && format == "" {
		return h.GetFile(c)
	}

	// Parse dimensions
	width := 0
	height := 0
	if widthStr != "" {
		if w, err := strconv.Atoi(widthStr); err == nil && w > 0 && w <= 3000 {
			width = w
		}
	}
	if heightStr != "" {
		if h, err := strconv.Atoi(heightStr); err == nil && h > 0 && h <= 3000 {
			height = h
		}
	}

	// Parse quality
	qualityInt := 85
	if quality != "" {
		if q, err := strconv.Atoi(quality); err == nil && q > 0 && q <= 100 {
			qualityInt = q
		}
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("%s_w%d_h%d_q%d", path, width, height, qualityInt)
	if format != "" {
		cacheKey += "_" + format
	}

	// TODO: Implement caching mechanism here

	// Open and resize image
	img, err := imaging.Open(imagePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "画像を開けませんでした",
		})
	}

	// Resize if dimensions are specified
	if width > 0 || height > 0 {
		if width > 0 && height > 0 {
			// Fit within bounds
			img = imaging.Fit(img, width, height, imaging.Lanczos)
		} else if width > 0 {
			// Resize by width
			img = imaging.Resize(img, width, 0, imaging.Lanczos)
		} else {
			// Resize by height
			img = imaging.Resize(img, 0, height, imaging.Lanczos)
		}
	}

	// Set appropriate content type
	contentType := GetMIMEType(imagePath)
	if format != "" {
		switch format {
		case "jpeg", "jpg":
			contentType = "image/jpeg"
		case "png":
			contentType = "image/png"
		case "webp":
			contentType = "image/webp"
		}
	}

	// Set headers
	c.Response().Header().Set("Content-Type", contentType)
	c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")

	// Encode and send image
	switch contentType {
	case "image/jpeg":
		return imaging.Encode(c.Response().Writer, img, imaging.JPEG, imaging.JPEGQuality(qualityInt))
	case "image/png":
		return imaging.Encode(c.Response().Writer, img, imaging.PNG)
	case "image/gif":
		return imaging.Encode(c.Response().Writer, img, imaging.GIF)
	default:
		// Default to JPEG
		return imaging.Encode(c.Response().Writer, img, imaging.JPEG, imaging.JPEGQuality(qualityInt))
	}
}