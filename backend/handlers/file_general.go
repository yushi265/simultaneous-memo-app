package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"simultaneous-memo-app/backend/models"
	"simultaneous-memo-app/backend/middleware"
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
)

const (
	MaxGeneralFileSize = 50 * 1024 * 1024 // 50MB
)

// AllowedFileTypes defines allowed file extensions and their MIME types
var AllowedFileTypes = map[string][]string{
	// Documents
	".pdf":  {"application/pdf"},
	".doc":  {"application/msword"},
	".docx": {"application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
	".xls":  {"application/vnd.ms-excel"},
	".xlsx": {"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
	".ppt":  {"application/vnd.ms-powerpoint"},
	".pptx": {"application/vnd.openxmlformats-officedocument.presentationml.presentation"},
	".txt":  {"text/plain"},
	".csv":  {"text/csv"},
	".rtf":  {"application/rtf"},
	
	// Archives
	".zip": {"application/zip"},
	".rar": {"application/x-rar-compressed"},
	".7z":  {"application/x-7z-compressed"},
	".tar": {"application/x-tar"},
	".gz":  {"application/gzip"},
	
	// Code files
	".js":   {"text/javascript", "application/javascript"},
	".ts":   {"text/typescript", "application/typescript"},
	".json": {"application/json"},
	".xml":  {"text/xml", "application/xml"},
	".html": {"text/html"},
	".css":  {"text/css"},
	".py":   {"text/x-python"},
	".go":   {"text/x-go"},
	".java": {"text/x-java"},
	".cpp":  {"text/x-c++"},
	".c":    {"text/x-c"},
	".sh":   {"text/x-shellscript"},
	".md":   {"text/markdown"},
}

// UploadGeneralFile handles non-image file uploads
func (h *Handler) UploadGeneralFile(c echo.Context) error {
	// Get workspace and user IDs from context
	workspaceID, err := middleware.GetWorkspaceID(c)
	if err != nil {
		return err
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file provided"})
	}

	// Validate file size
	if file.Size > MaxGeneralFileSize {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("File size exceeds %dMB limit", MaxGeneralFileSize/(1024*1024))})
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	
	// Check if file type is allowed
	allowedMimes, isAllowed := AllowedFileTypes[ext]
	if !isAllowed {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File type not allowed"})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
	}
	defer src.Close()

	// Detect content type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read file"})
	}
	contentType := http.DetectContentType(buffer)
	src.Seek(0, 0) // Reset file pointer

	// Validate content type matches allowed types
	validType := false
	for _, allowed := range allowedMimes {
		if strings.Contains(contentType, allowed) || contentType == "application/octet-stream" {
			validType = true
			break
		}
	}
	if !validType {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid file content type"})
	}

	// Create directory structure
	uploadDir := filepath.Join("uploads", "files", time.Now().Format("2006/01"))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	safeFilename := sanitizeGeneralFilename(file.Filename)
	filename := fmt.Sprintf("%d_%s", timestamp, safeFilename)
	filePath := filepath.Join(uploadDir, filename)

	// Save the file
	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create file"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		os.Remove(filePath)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}

	// Get page ID if provided
	var pageID *uuid.UUID
	if pageIDStr := c.FormValue("page_id"); pageIDStr != "" {
		if id, err := uuid.Parse(pageIDStr); err == nil {
			pageID = &id
		}
	}

	// Save file metadata to database
	fileModel := &models.File{
		WorkspaceID:  workspaceID,
		UserID:       userID,
		Filename:     filename,
		OriginalName: file.Filename,
		ContentType:  contentType,
		Size:         file.Size,
		Path:         filePath,
		PageID:       pageID,
	}

	if err := h.db.Create(fileModel).Error; err != nil {
		os.Remove(filePath)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file metadata"})
	}

	// Return file metadata
	return c.JSON(http.StatusOK, fileModel.ToMetadata(getBaseURL(c)))
}

// ListFiles returns files with pagination and optional filtering
func (h *Handler) ListFiles(c echo.Context) error {
	// Get workspace ID from context
	workspaceID, err := middleware.GetWorkspaceID(c)
	if err != nil {
		return err
	}

	var files []models.File
	query := h.db.Model(&models.File{}).Where("workspace_id = ?", workspaceID)

	// Pagination parameters
	page := 1
	limit := 20
	if p := c.QueryParam("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	if l := c.QueryParam("limit"); l != "" {
		if limitNum, err := strconv.Atoi(l); err == nil && limitNum > 0 && limitNum <= 100 {
			limit = limitNum
		}
	}
	offset := (page - 1) * limit

	// Count total files
	var total int64
	countQuery := h.db.Model(&models.File{}).Where("workspace_id = ?", workspaceID)

	// Filter by page ID if provided
	if pageID := c.QueryParam("page_id"); pageID != "" {
		query = query.Where("page_id = ?", pageID)
		countQuery = countQuery.Where("page_id = ?", pageID)
	}

	// Filter by file type if provided
	if fileType := c.QueryParam("type"); fileType != "" {
		switch fileType {
		case "document":
			query = query.Where("content_type LIKE ?", "%document%").Or("content_type LIKE ?", "%pdf%").Or("content_type LIKE ?", "%text%")
			countQuery = countQuery.Where("content_type LIKE ?", "%document%").Or("content_type LIKE ?", "%pdf%").Or("content_type LIKE ?", "%text%")
		case "archive":
			query = query.Where("content_type LIKE ?", "%zip%").Or("content_type LIKE ?", "%compressed%")
			countQuery = countQuery.Where("content_type LIKE ?", "%zip%").Or("content_type LIKE ?", "%compressed%")
		case "code":
			query = query.Where("content_type LIKE ?", "%javascript%").Or("content_type LIKE ?", "%json%").Or("content_type LIKE ?", "%xml%")
			countQuery = countQuery.Where("content_type LIKE ?", "%javascript%").Or("content_type LIKE ?", "%json%").Or("content_type LIKE ?", "%xml%")
		}
	}

	// Get total count
	if err := countQuery.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to count files"})
	}

	// Get paginated files
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&files).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch files"})
	}

	// Convert to metadata
	baseURL := getBaseURL(c)
	metadata := make([]models.FileMetadata, len(files))
	for i, file := range files {
		metadata[i] = file.ToMetadata(baseURL)
	}

	// Calculate pagination info
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"files":       metadata,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"has_more":    page < totalPages,
	})
}

// GetFileMetadata returns metadata for a specific file
func (h *Handler) GetFileMetadata(c echo.Context) error {
	// Get workspace ID from context
	workspaceID, err := middleware.GetWorkspaceID(c)
	if err != nil {
		return err
	}

	id := c.Param("id")
	
	var file models.File
	if err := h.db.Where("id = ? AND workspace_id = ?", id, workspaceID).First(&file).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "File not found"})
	}

	return c.JSON(http.StatusOK, file.ToMetadata(getBaseURL(c)))
}

// DeleteFile deletes a file and its metadata
func (h *Handler) DeleteFile(c echo.Context) error {
	// Get workspace ID from context
	workspaceID, err := middleware.GetWorkspaceID(c)
	if err != nil {
		return err
	}

	id := c.Param("id")
	
	var file models.File
	if err := h.db.Where("id = ? AND workspace_id = ?", id, workspaceID).First(&file).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "File not found"})
	}

	// Delete physical file
	if err := os.Remove(file.Path); err != nil && !os.IsNotExist(err) {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete file"})
	}

	// Delete database record
	if err := h.db.Delete(&file).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete file metadata"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "File deleted successfully"})
}

// ServeFile serves uploaded files
func (h *Handler) ServeFile(c echo.Context) error {
	// Get workspace ID from context
	workspaceID, err := middleware.GetWorkspaceID(c)
	if err != nil {
		return err
	}

	filename := c.Param("*")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No filename provided"})
	}

	// Validate filename - no path traversal allowed
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid filename"})
	}

	// Find file in database
	var file models.File
	if err := h.db.Where("filename = ? AND workspace_id = ?", filename, workspaceID).First(&file).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "File not found"})
	}

	// Validate that the file path is within uploads directory
	absPath, err := filepath.Abs(file.Path)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid file path"})
	}
	
	uploadsDir, err := filepath.Abs("uploads")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server configuration error"})
	}
	
	// Ensure the file path is within the uploads directory
	if !strings.HasPrefix(absPath, uploadsDir) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	// Check if file exists
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "File not found on disk"})
	}

	// Set appropriate headers
	c.Response().Header().Set("Content-Type", file.ContentType)
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))
	
	// Set cache headers for static files
	c.Response().Header().Set("Cache-Control", "public, max-age=86400") // 1 day

	// Determine if file should be downloaded or displayed inline
	disposition := "inline"
	if c.QueryParam("download") == "true" {
		disposition = "attachment"
	}
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`%s; filename="%s"`, disposition, file.OriginalName))

	return c.File(file.Path)
}

// sanitizeGeneralFilename removes potentially dangerous characters from filename
func sanitizeGeneralFilename(filename string) string {
	// Remove path separators and other dangerous characters
	filename = filepath.Base(filename)
	filename = strings.ReplaceAll(filename, "..", "")
	
	// Replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")
	
	// Keep the file extension
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	
	// For the name part, replace non-ASCII characters with underscores
	// but keep basic alphanumeric, dots, dashes, and underscores
	safeName := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, nameWithoutExt)
	
	// Remove multiple consecutive underscores
	for strings.Contains(safeName, "__") {
		safeName = strings.ReplaceAll(safeName, "__", "_")
	}
	
	// Trim underscores from start and end
	safeName = strings.Trim(safeName, "_")
	
	// If the name is empty after sanitization, use a default
	if safeName == "" {
		safeName = "file"
	}
	
	return safeName + ext
}

// getBaseURL constructs the base URL from the request
func getBaseURL(c echo.Context) string {
	scheme := "http"
	if c.Request().TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request().Host)
}