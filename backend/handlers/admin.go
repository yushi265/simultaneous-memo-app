package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// CleanupImages removes orphaned images older than 24 hours
func (h *Handler) CleanupImages(c echo.Context) error {
	// This endpoint should be protected in production
	// For now, we'll add a simple token check
	token := c.QueryParam("token")
	if token != "cleanup-secret-token" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "認証エラー",
		})
	}

	// Cleanup orphaned images older than 24 hours
	err := CleanupOrphanedImages(h.db, 24*time.Hour)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "クリーンアップ中にエラーが発生しました",
			"details": err.Error(),
		})
	}

	// Remove empty directories
	if err := RemoveEmptyDirectories("../uploads/images"); err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "画像のクリーンアップは完了しましたが、ディレクトリ削除でエラーがありました",
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "クリーンアップが正常に完了しました",
	})
}