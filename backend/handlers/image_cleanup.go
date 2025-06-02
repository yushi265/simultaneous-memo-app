package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"simultaneous-memo-app/backend/models"
	"gorm.io/gorm"
)

// DeleteImagesByPageID deletes all images associated with a page
func DeleteImagesByPageID(db *gorm.DB, pageID uint) error {
	// Get all images for the page
	images, err := models.GetImagesByPageID(db, pageID)
	if err != nil {
		return fmt.Errorf("画像の取得に失敗しました: %w", err)
	}

	errors := []error{}

	// Delete each image
	for _, image := range images {
		// Delete main image file
		mainPath := filepath.Join("../uploads", image.Path)
		if err := os.Remove(mainPath); err != nil && !os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("画像 %s の削除エラー: %w", image.Filename, err))
		}

		// Delete thumbnail
		if image.ThumbnailPath != "" {
			thumbPath := filepath.Join("../uploads", image.ThumbnailPath)
			if err := os.Remove(thumbPath); err != nil && !os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf("サムネイル %s の削除エラー: %w", image.Filename, err))
			}
		}

		// Delete from database
		if err := models.DeleteImage(db, image.ID); err != nil {
			errors = append(errors, fmt.Errorf("画像 %s のDB削除エラー: %w", image.Filename, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("画像削除中にエラーが発生しました: %v", errors)
	}

	return nil
}

// CleanupOrphanedImages removes images not associated with any page
func CleanupOrphanedImages(db *gorm.DB, olderThan time.Duration) error {
	// Get orphaned images older than specified duration
	cutoffTime := time.Now().Add(-olderThan)
	orphanedImages, err := models.GetOrphanedImages(db, cutoffTime)
	if err != nil {
		return fmt.Errorf("孤立した画像の取得に失敗しました: %w", err)
	}

	deletedCount := 0
	errors := []error{}

	for _, image := range orphanedImages {
		// Delete main image file
		mainPath := filepath.Join("../uploads", image.Path)
		if err := os.Remove(mainPath); err != nil && !os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("画像 %s の削除エラー: %w", image.Filename, err))
			continue
		}

		// Delete thumbnail
		if image.ThumbnailPath != "" {
			thumbPath := filepath.Join("../uploads", image.ThumbnailPath)
			if err := os.Remove(thumbPath); err != nil && !os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf("サムネイル %s の削除エラー: %w", image.Filename, err))
			}
		}

		// Delete from database
		if err := models.DeleteImage(db, image.ID); err != nil {
			errors = append(errors, fmt.Errorf("画像 %s のDB削除エラー: %w", image.Filename, err))
			continue
		}

		deletedCount++
	}

	if len(errors) > 0 {
		return fmt.Errorf("%d個の画像を削除しましたが、エラーがありました: %v", deletedCount, errors)
	}

	return nil
}

// RemoveEmptyDirectories removes empty upload directories
func RemoveEmptyDirectories(basePath string) error {
	return filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a directory
		if !info.IsDir() {
			return nil
		}

		// Skip base directory
		if path == basePath {
			return nil
		}

		// Check if directory is empty
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		// Remove if empty
		if len(entries) == 0 {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("空のディレクトリ %s の削除エラー: %w", path, err)
			}
		}

		return nil
	})
}