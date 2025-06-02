package handlers

import (
	"path/filepath"
	"strings"
)

// ExtensionToMIME maps file extensions to MIME types
var ExtensionToMIME = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
	".svg":  "image/svg+xml",
	".bmp":  "image/bmp",
	".ico":  "image/x-icon",
	".tiff": "image/tiff",
	".tif":  "image/tiff",
}

// GetMIMEType returns the MIME type for a file based on its extension
func GetMIMEType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if mimeType, ok := ExtensionToMIME[ext]; ok {
		return mimeType
	}
	return "application/octet-stream"
}

// IsImageMIME checks if a MIME type is an image
func IsImageMIME(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}

// GetImageHeaders returns appropriate headers for image responses
func GetImageHeaders(filename string, fileSize int64) map[string]string {
	headers := make(map[string]string)
	
	// Set MIME type
	headers["Content-Type"] = GetMIMEType(filename)
	
	// Set security headers
	headers["X-Content-Type-Options"] = "nosniff"
	
	// Set performance headers
	headers["Accept-Ranges"] = "bytes"
	
	// For SVG files, add CSP header for security
	if strings.HasSuffix(strings.ToLower(filename), ".svg") {
		headers["Content-Security-Policy"] = "default-src 'none'; style-src 'unsafe-inline'; sandbox"
	}
	
	return headers
}