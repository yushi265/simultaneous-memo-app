package handlers

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// ImageConfig holds image processing configuration
type ImageConfig struct {
	MaxWidth       int
	MaxHeight      int
	ThumbnailWidth int
	Quality        int
}

// DefaultImageConfig returns default image configuration
func DefaultImageConfig() ImageConfig {
	return ImageConfig{
		MaxWidth:       1920,
		MaxHeight:      1080,
		ThumbnailWidth: 300,
		Quality:        85,
	}
}

// ProcessImage resizes and optimizes an image
func ProcessImage(srcPath, dstPath string, config ImageConfig) error {
	// Open the source image with auto-orientation
	src, err := imaging.Open(srcPath, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("画像を開けませんでした: %w", err)
	}

	// Get original dimensions
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate new dimensions while maintaining aspect ratio
	if width > config.MaxWidth || height > config.MaxHeight {
		src = imaging.Fit(src, config.MaxWidth, config.MaxHeight, imaging.Lanczos)
	}

	// Determine output format based on file extension
	ext := filepath.Ext(dstPath)
	switch ext {
	case ".png":
		err = imaging.Save(src, dstPath, imaging.PNGCompressionLevel(6))
	case ".gif":
		err = imaging.Save(src, dstPath)
	case ".webp":
		// WebP is not directly supported by imaging library, save as JPEG
		dstPath = dstPath[:len(dstPath)-5] + ".jpg"
		err = imaging.Save(src, dstPath, imaging.JPEGQuality(config.Quality))
	default:
		// Default to JPEG
		err = imaging.Save(src, dstPath, imaging.JPEGQuality(config.Quality))
	}

	if err != nil {
		return fmt.Errorf("画像の保存に失敗しました: %w", err)
	}

	return nil
}

// CreateThumbnail creates a thumbnail version of the image
func CreateThumbnail(srcPath, thumbPath string, config ImageConfig) error {
	// Open the source image
	src, err := imaging.Open(srcPath)
	if err != nil {
		return fmt.Errorf("画像を開けませんでした: %w", err)
	}

	// Create thumbnail with fixed width
	thumbnail := imaging.Resize(src, config.ThumbnailWidth, 0, imaging.Lanczos)

	// Ensure thumbnail directory exists
	thumbDir := filepath.Dir(thumbPath)
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return fmt.Errorf("サムネイルディレクトリの作成に失敗しました: %w", err)
	}

	// Save thumbnail
	err = imaging.Save(thumbnail, thumbPath, imaging.JPEGQuality(config.Quality))
	if err != nil {
		return fmt.Errorf("サムネイルの保存に失敗しました: %w", err)
	}

	return nil
}

// GetImageDimensions returns the dimensions of an image
func GetImageDimensions(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}