package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
	
	"simultaneous-memo-app/backend/models"
	"simultaneous-memo-app/backend/middleware"
	"simultaneous-memo-app/backend/auth"
)

type WorkspaceHandler struct {
	db *gorm.DB
}

func NewWorkspaceHandler(db *gorm.DB) *WorkspaceHandler {
	return &WorkspaceHandler{db: db}
}

// CreateWorkspaceRequest represents workspace creation request
type CreateWorkspaceRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// UpdateWorkspaceRequest represents workspace update request
type UpdateWorkspaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// WorkspaceResponse represents workspace with member info
type WorkspaceResponse struct {
	models.Workspace
	MemberCount int    `json:"member_count"`
	UserRole    string `json:"user_role"`
}

// GetWorkspaces retrieves all workspaces the user is a member of
func (h *WorkspaceHandler) GetWorkspaces(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	var memberships []models.WorkspaceMember
	if err := h.db.Preload("Workspace").Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get workspaces")
	}

	workspaces := make([]WorkspaceResponse, len(memberships))
	for i, membership := range memberships {
		// Get member count for each workspace
		var memberCount int64
		h.db.Model(&models.WorkspaceMember{}).Where("workspace_id = ?", membership.WorkspaceID).Count(&memberCount)

		workspaces[i] = WorkspaceResponse{
			Workspace:   membership.Workspace,
			MemberCount: int(memberCount),
			UserRole:    membership.Role,
		}
	}

	return c.JSON(http.StatusOK, workspaces)
}

// GetWorkspace retrieves a specific workspace
func (h *WorkspaceHandler) GetWorkspace(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workspace ID")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	// Check if user is a member of this workspace
	var membership models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&membership).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied to this workspace")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	// Get workspace details
	var workspace models.Workspace
	if err := h.db.Preload("Owner").First(&workspace, workspaceID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Workspace not found")
	}

	// Get member count
	var memberCount int64
	h.db.Model(&models.WorkspaceMember{}).Where("workspace_id = ?", workspaceID).Count(&memberCount)

	response := WorkspaceResponse{
		Workspace:   workspace,
		MemberCount: int(memberCount),
		UserRole:    membership.Role,
	}

	return c.JSON(http.StatusOK, response)
}

// CreateWorkspace creates a new workspace
func (h *WorkspaceHandler) CreateWorkspace(c echo.Context) error {
	var req CreateWorkspaceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	// Generate unique slug from name
	slug := generateSlug(req.Name)
	
	// Check if slug is unique
	for {
		var existingWorkspace models.Workspace
		if err := h.db.Where("slug = ?", slug).First(&existingWorkspace).Error; err == gorm.ErrRecordNotFound {
			break // Slug is unique
		}
		// Add random suffix if slug exists
		slug = generateSlug(req.Name) + "-" + uuid.New().String()[:8]
	}

	// Create workspace
	workspace := models.Workspace{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		IsPersonal:  false,
		OwnerID:     userID,
	}

	if err := h.db.Create(&workspace).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create workspace")
	}

	// Add creator as owner member
	member := models.WorkspaceMember{
		WorkspaceID: workspace.ID,
		UserID:      userID,
		Role:        "owner",
	}

	if err := h.db.Create(&member).Error; err != nil {
		// Rollback workspace creation
		h.db.Delete(&workspace)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add workspace member")
	}

	response := WorkspaceResponse{
		Workspace:   workspace,
		MemberCount: 1,
		UserRole:    "owner",
	}

	return c.JSON(http.StatusCreated, response)
}

// UpdateWorkspace updates an existing workspace
func (h *WorkspaceHandler) UpdateWorkspace(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workspace ID")
	}

	var req UpdateWorkspaceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	// Check if user has admin or owner role
	var membership models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&membership).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied to this workspace")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	if membership.Role != "owner" && membership.Role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to update workspace")
	}

	// Get workspace
	var workspace models.Workspace
	if err := h.db.First(&workspace, workspaceID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Workspace not found")
	}

	// Prevent editing personal workspaces
	if workspace.IsPersonal {
		return echo.NewHTTPError(http.StatusForbidden, "Cannot edit personal workspace")
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
		// Update slug if name changed
		if req.Name != workspace.Name {
			updates["slug"] = generateSlug(req.Name)
		}
	}
	if req.Description != workspace.Description {
		updates["description"] = req.Description
	}

	if len(updates) > 0 {
		if err := h.db.Model(&workspace).Updates(updates).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update workspace")
		}
	}

	// Get updated workspace
	if err := h.db.First(&workspace, workspaceID).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get updated workspace")
	}

	// Get member count
	var memberCount int64
	h.db.Model(&models.WorkspaceMember{}).Where("workspace_id = ?", workspaceID).Count(&memberCount)

	response := WorkspaceResponse{
		Workspace:   workspace,
		MemberCount: int(memberCount),
		UserRole:    membership.Role,
	}

	return c.JSON(http.StatusOK, response)
}

// DeleteWorkspace deletes a workspace (only owner can delete)
func (h *WorkspaceHandler) DeleteWorkspace(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workspace ID")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	// Get workspace
	var workspace models.Workspace
	if err := h.db.First(&workspace, workspaceID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Workspace not found")
	}

	// Prevent deleting personal workspaces
	if workspace.IsPersonal {
		return echo.NewHTTPError(http.StatusForbidden, "Cannot delete personal workspace")
	}

	// Check if user is the owner
	if workspace.OwnerID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Only workspace owner can delete workspace")
	}

	// Delete workspace (cascading deletes will handle members, pages, etc.)
	if err := h.db.Delete(&workspace).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete workspace")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Workspace deleted successfully",
	})
}

// SwitchWorkspace switches user's current workspace
func (h *WorkspaceHandler) SwitchWorkspace(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workspace ID")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	// Check if user is a member of this workspace
	var membership models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&membership).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied to this workspace")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	// Get workspace details
	var workspace models.Workspace
	if err := h.db.First(&workspace, workspaceID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Workspace not found")
	}

	// Get user details
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Generate new token with updated workspace
	newToken, err := generateTokenWithWorkspace(user, workspace)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":     newToken,
		"workspace": workspace,
		"message":   "Workspace switched successfully",
	})
}

// Helper function to generate slug from name
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Helper function to generate token with workspace
func generateTokenWithWorkspace(user models.User, workspace models.Workspace) (string, error) {
	return auth.GenerateToken(user.ID, user.Email, user.Name, workspace.ID)
}