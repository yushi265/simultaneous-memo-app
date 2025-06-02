package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"simultaneous-memo-app/backend/models"

	"github.com/labstack/echo/v4"
)

// GetImages retrieves all images or images for a specific page
func (h *Handler) GetImages(c echo.Context) error {
	pageIDStr := c.QueryParam("page_id")
	
	if pageIDStr != "" {
		// Get images for specific page
		pageID, err := strconv.ParseUint(pageIDStr, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "無効なページIDです",
			})
		}

		images, err := models.GetImagesByPageID(h.db, uint(pageID))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "画像の取得に失敗しました",
			})
		}

		return c.JSON(http.StatusOK, images)
	}

	// Get all images
	var images []models.Image
	if err := h.db.Order("created_at DESC").Find(&images).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "画像の取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, images)
}

// GetImageByID retrieves a single image by ID
func (h *Handler) GetImageByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効な画像IDです",
		})
	}

	image, err := models.GetImageByID(h.db, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "画像が見つかりません",
		})
	}

	return c.JSON(http.StatusOK, image)
}

// DeleteImageByID deletes an image and its files
func (h *Handler) DeleteImageByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効な画像IDです",
		})
	}

	// Get image metadata
	image, err := models.GetImageByID(h.db, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "画像が見つかりません",
		})
	}

	// Delete actual files from filesystem
	errors := []string{}
	
	// Delete main image file
	mainPath := filepath.Join("../uploads", image.Path)
	if err := os.Remove(mainPath); err != nil && !os.IsNotExist(err) {
		errors = append(errors, fmt.Sprintf("メイン画像の削除エラー: %v", err))
	}
	
	// Delete thumbnail file
	if image.ThumbnailPath != "" {
		thumbPath := filepath.Join("../uploads", image.ThumbnailPath)
		if err := os.Remove(thumbPath); err != nil && !os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("サムネイルの削除エラー: %v", err))
		}
	}

	// Delete from database
	if err := models.DeleteImage(h.db, uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "データベースからの削除に失敗しました",
		})
	}

	// Return response
	if len(errors) > 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "画像を削除しましたが、一部エラーがありました",
			"errors":  errors,
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "画像を正常に削除しました",
	})
}