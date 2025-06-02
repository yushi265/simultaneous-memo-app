package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"simultaneous-memo-app/backend/models"

	"github.com/labstack/echo/v4"
)

const (
	// MaxFileSize is 10MB
	MaxFileSize = 10 * 1024 * 1024
)

// AllowedImageTypes contains allowed MIME types for images
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// AllowedImageExtensions contains allowed file extensions
var AllowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// UploadFile handles file uploads with validation and saves metadata to database
func (h *Handler) UploadFile(c echo.Context) error {
	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ファイルがアップロードされていません",
		})
	}

	// Validate file size
	if file.Size > MaxFileSize {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("ファイルサイズが大きすぎます。最大サイズは%dMBです", MaxFileSize/1024/1024),
		})
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedImageExtensions[ext] {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "許可されていないファイル形式です。JPEG、PNG、GIF、WebPのみアップロード可能です",
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ファイルを開けませんでした",
		})
	}
	defer src.Close()

	// Read first 512 bytes to detect MIME type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ファイルの読み取りに失敗しました",
		})
	}

	// Detect content type
	contentType := http.DetectContentType(buffer)
	if !AllowedImageTypes[contentType] {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "許可されていないファイル形式です。画像ファイルのみアップロード可能です",
		})
	}

	// Reset file reader to beginning
	src.Seek(0, 0)

	// Create uploads directory structure (YYYY/MM)
	now := time.Now()
	uploadsDir := fmt.Sprintf("../uploads/images/%d/%02d", now.Year(), now.Month())
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "アップロードディレクトリの作成に失敗しました",
		})
	}

	// Generate unique filename with sanitization
	timestamp := now.Unix()
	safeFilename := sanitizeFilename(file.Filename)
	filename := fmt.Sprintf("%d_%s", timestamp, safeFilename)
	filepath := filepath.Join(uploadsDir, filename)

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ファイルの作成に失敗しました",
		})
	}
	defer dst.Close()

	// Copy file first to temporary location
	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ファイルの保存に失敗しました",
		})
	}
	dst.Close()

	// Process image (resize and optimize)
	config := DefaultImageConfig()
	processedPath := filepath
	
	// For JPEG and PNG, apply processing
	if contentType == "image/jpeg" || contentType == "image/png" {
		err = ProcessImage(filepath, processedPath, config)
		if err != nil {
			// If processing fails, keep the original
			processedPath = filepath
		}
	}

	// Create thumbnail
	thumbFilename := "thumb_" + filename
	thumbPath := filepath.Join(uploadsDir, thumbFilename)
	err = CreateThumbnail(processedPath, thumbPath, config)
	if err != nil {
		// Log error but don't fail the upload
		fmt.Printf("サムネイル作成エラー: %v\n", err)
	}

	// Get image dimensions
	width, height, err := GetImageDimensions(processedPath)
	if err != nil {
		width, height = 0, 0
	}

	// Get final file size after processing
	fileInfo, _ := os.Stat(processedPath)
	finalSize := file.Size
	if fileInfo != nil {
		finalSize = fileInfo.Size()
	}

	// Save image metadata to database
	relativePath := fmt.Sprintf("/images/%d/%02d/%s", now.Year(), now.Month(), filename)
	thumbRelativePath := fmt.Sprintf("/images/%d/%02d/%s", now.Year(), now.Month(), thumbFilename)
	
	imageRecord := &models.Image{
		Filename:      filename,
		OriginalName:  file.Filename,
		Path:          relativePath,
		ThumbnailPath: thumbRelativePath,
		Size:          finalSize,
		Width:         width,
		Height:        height,
		ContentType:   contentType,
	}

	// Check if page_id is provided in the request
	if pageIDStr := c.FormValue("page_id"); pageIDStr != "" {
		if pageID, err := strconv.ParseUint(pageIDStr, 10, 32); err == nil {
			pageIDUint := uint(pageID)
			imageRecord.PageID = &pageIDUint
		}
	}

	// Save to database
	if err := models.CreateImage(h.db, imageRecord); err != nil {
		// Log error but don't fail the upload
		fmt.Printf("画像メタデータの保存エラー: %v\n", err)
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":          imageRecord.ID,
		"filename":    filename,
		"size":        finalSize,
		"originalSize": file.Size,
		"url":         fmt.Sprintf("/api/files%s", relativePath),
		"thumbnailUrl": fmt.Sprintf("/api/files%s", thumbRelativePath),
		"contentType": contentType,
		"width":       width,
		"height":      height,
		"uploadedAt":  imageRecord.CreatedAt,
	})
}

// sanitizeFilename removes potentially dangerous characters from filename
func sanitizeFilename(filename string) string {
	// Get file extension
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	
	// Remove or replace dangerous characters
	reg := regexp.MustCompile(`[^a-zA-Z0-9\-_]`)
	safeName := reg.ReplaceAllString(name, "_")
	
	// Ensure filename is not empty
	if safeName == "" {
		safeName = "file"
	}
	
	// Limit filename length
	if len(safeName) > 100 {
		safeName = safeName[:100]
	}
	
	return safeName + ext
}

// GetFile serves uploaded files with proper MIME types and caching
func (h *Handler) GetFile(c echo.Context) error {
	// Parse the path parameter to support nested directories
	path := c.Param("*")
	filepath := filepath.Join("../uploads", path)

	// Prevent path traversal attacks
	if strings.Contains(path, "..") {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "不正なパスです",
		})
	}

	// Check if file exists
	fileInfo, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "ファイルが見つかりません",
		})
	}

	// Open file to detect MIME type
	file, err := os.Open(filepath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ファイルを開けませんでした",
		})
	}
	defer file.Close()

	// Get MIME type based on file extension first
	contentType := GetMIMEType(filepath)
	
	// If unknown, detect from file content
	if contentType == "application/octet-stream" {
		buffer := make([]byte, 512)
		n, _ := file.Read(buffer)
		contentType = http.DetectContentType(buffer[:n])
		// Reset file pointer
		file.Seek(0, 0)
	}

	// Get appropriate headers for the file
	headers := GetImageHeaders(filepath, fileInfo.Size())
	for key, value := range headers {
		c.Response().Header().Set(key, value)
	}
	
	// Set cache headers for images (1 year)
	if strings.HasPrefix(contentType, "image/") {
		c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		c.Response().Header().Set("Expires", time.Now().Add(365*24*time.Hour).Format(http.TimeFormat))
	}

	// Set Content-Disposition for download
	filename := filepath.Base(filepath)
	if c.QueryParam("download") == "true" {
		c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	} else {
		c.Response().Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	}

	// Set Last-Modified header
	c.Response().Header().Set("Last-Modified", fileInfo.ModTime().Format(http.TimeFormat))

	// Handle If-Modified-Since header for caching
	if ifModifiedSince := c.Request().Header.Get("If-Modified-Since"); ifModifiedSince != "" {
		t, err := time.Parse(http.TimeFormat, ifModifiedSince)
		if err == nil && !fileInfo.ModTime().After(t) {
			return c.NoContent(http.StatusNotModified)
		}
	}

	// Serve the file
	return c.Stream(http.StatusOK, contentType, file)
}