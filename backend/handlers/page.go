package handlers

import (
	"fmt"
	"net/http"

	"simultaneous-memo-app/backend/models"
	"simultaneous-memo-app/backend/middleware"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GetPages retrieves all pages in the user's workspace
func (h *Handler) GetPages(c echo.Context) error {
	workspaceID, err := middleware.GetWorkspaceID(c)
	if err != nil {
		return err
	}

	pages, err := models.GetAllPages(h.db, workspaceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve pages",
		})
	}

	return c.JSON(http.StatusOK, pages)
}

// GetPage retrieves a single page by ID
func (h *Handler) GetPage(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid page ID",
		})
	}

	workspaceID, err := middleware.GetWorkspaceID(c)
	if err != nil {
		return err
	}

	page, err := models.GetPageByID(h.db, id, workspaceID)
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

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	workspaceID, err := middleware.GetWorkspaceID(c)
	if err != nil {
		return err
	}

	// Set workspace ID
	page.WorkspaceID = workspaceID

	// Set default content if not provided
	if page.Content == nil {
		page.Content = []byte(`{"doc":{"type":"doc","content":[]}}`)
	}

	if err := models.CreatePage(h.db, &page, userID); err != nil {
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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid page ID",
		})
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	workspaceID, err := middleware.GetWorkspaceID(c)
	if err != nil {
		return err
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := models.UpdatePage(h.db, id, workspaceID, updates, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update page",
		})
	}

	// Update image references if content was updated
	if content, ok := updates["content"]; ok {
		if contentJSON, ok := content.([]byte); ok {
			if err := models.UpdateImageReferences(h.db, id, contentJSON); err != nil {
				// Log error but don't fail the request
				fmt.Printf("画像参照の更新エラー: %v\n", err)
			}
		}
	}

	page, err := models.GetPageByID(h.db, id, workspaceID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Page not found",
		})
	}

	return c.JSON(http.StatusOK, page)
}

// DeletePage deletes a page and its associated images
func (h *Handler) DeletePage(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なページIDです",
		})
	}

	workspaceID, err := middleware.GetWorkspaceID(c)
	if err != nil {
		return err
	}

	// Delete associated images first
	if err := DeleteImagesByPageID(h.db, id, workspaceID); err != nil {
		// Log error but continue with page deletion
		fmt.Printf("ページ %v の画像削除エラー: %v\n", id, err)
	}

	// Delete the page
	if err := models.DeletePage(h.db, id, workspaceID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ページの削除に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "ページと関連画像を削除しました",
	})
}