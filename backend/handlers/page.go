package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"simultaneous-memo-app/backend/models"

	"github.com/labstack/echo/v4"
)

// GetPages retrieves all pages
func (h *Handler) GetPages(c echo.Context) error {
	pages, err := models.GetAllPages(h.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve pages",
		})
	}

	return c.JSON(http.StatusOK, pages)
}

// GetPage retrieves a single page by ID
func (h *Handler) GetPage(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid page ID",
		})
	}

	page, err := models.GetPageByID(h.db, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Page not found",
		})
	}

	return c.JSON(http.StatusOK, page)
}

// CreatePage creates a new page
func (h *Handler) CreatePage(c echo.Context) error {
	var page models.Page
	if err := c.Bind(&page); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Set default content if not provided
	if page.Content == nil {
		page.Content = []byte(`{"doc":{"type":"doc","content":[]}}`)
	}

	if err := models.CreatePage(h.db, &page); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create page",
		})
	}

	// Update image references
	if err := models.UpdateImageReferences(h.db, page.ID, page.Content); err != nil {
		// Log error but don't fail the request
		fmt.Printf("画像参照の更新エラー: %v\n", err)
	}

	return c.JSON(http.StatusCreated, page)
}

// UpdatePage updates an existing page
func (h *Handler) UpdatePage(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid page ID",
		})
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := models.UpdatePage(h.db, uint(id), updates); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update page",
		})
	}

	// Update image references if content was updated
	if content, ok := updates["content"]; ok {
		if contentJSON, ok := content.([]byte); ok {
			if err := models.UpdateImageReferences(h.db, uint(id), contentJSON); err != nil {
				// Log error but don't fail the request
				fmt.Printf("画像参照の更新エラー: %v\n", err)
			}
		}
	}

	page, err := models.GetPageByID(h.db, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Page not found",
		})
	}

	return c.JSON(http.StatusOK, page)
}

// DeletePage deletes a page and its associated images
func (h *Handler) DeletePage(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なページIDです",
		})
	}

	// Delete associated images first
	if err := DeleteImagesByPageID(h.db, uint(id)); err != nil {
		// Log error but continue with page deletion
		fmt.Printf("ページ %d の画像削除エラー: %v\n", id, err)
	}

	// Delete the page
	if err := models.DeletePage(h.db, uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ページの削除に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "ページと関連画像を削除しました",
	})
}