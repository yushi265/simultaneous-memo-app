package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

// InviteMemberRequest represents member invitation request
type InviteMemberRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin member viewer"`
}

// InviteMember creates an invitation for a user to join the workspace
func (h *WorkspaceHandler) InviteMember(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workspace ID")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	var req InviteMemberRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check if user has permission to invite (owner or admin)
	var membership models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&membership).Error; err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	if membership.Role != "owner" && membership.Role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to invite members")
	}

	// Check if workspace exists and is not personal
	var workspace models.Workspace
	if err := h.db.First(&workspace, workspaceID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Workspace not found")
	}

	if workspace.IsPersonal {
		return echo.NewHTTPError(http.StatusForbidden, "Cannot invite members to personal workspace")
	}

	// Check if user is already a member
	var existingMember models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id IN (SELECT id FROM users WHERE email = ?)", workspaceID, req.Email).First(&existingMember).Error; err == nil {
		return echo.NewHTTPError(http.StatusConflict, "User is already a member of this workspace")
	}

	// Check if there's already a pending invitation
	var existingInvitation models.WorkspaceInvitation
	if err := h.db.Where("workspace_id = ? AND email = ? AND accepted_at IS NULL AND expires_at > ?", workspaceID, req.Email, time.Now()).First(&existingInvitation).Error; err == nil {
		return echo.NewHTTPError(http.StatusConflict, "Invitation already sent to this email")
	}

	// Generate invitation token
	token, err := auth.GenerateInvitationToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate invitation token")
	}

	// Create invitation
	invitation := models.WorkspaceInvitation{
		WorkspaceID: workspaceID,
		InviterID:   userID,
		Email:       req.Email,
		Role:        req.Role,
		Token:       token,
		ExpiresAt:   auth.GetInvitationExpiration(),
	}

	if err := h.db.Create(&invitation).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create invitation")
	}

	// Load invitation with relationships for response
	if err := h.db.Preload("Workspace").Preload("Inviter").First(&invitation, invitation.ID).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load invitation details")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"invitation": invitation,
		"invite_url": fmt.Sprintf("/invite/%s", token),
		"message":    "Invitation sent successfully",
	})
}

// AcceptInvitation accepts a workspace invitation
func (h *WorkspaceHandler) AcceptInvitation(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid invitation token")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	// Get user email for validation
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user information")
	}

	// Find invitation
	var invitation models.WorkspaceInvitation
	if err := h.db.Preload("Workspace").Where("token = ?", token).First(&invitation).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid or expired invitation")
	}

	// Check if invitation is for the authenticated user
	if invitation.Email != user.Email {
		return echo.NewHTTPError(http.StatusForbidden, "This invitation is not for your email address")
	}

	// Check if already accepted
	if invitation.AcceptedAt != nil {
		return echo.NewHTTPError(http.StatusConflict, "Invitation already accepted")
	}

	// Check if expired
	if auth.IsInvitationExpired(invitation.ExpiresAt) {
		return echo.NewHTTPError(http.StatusGone, "Invitation has expired")
	}

	// Check if user is already a member
	var existingMember models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id = ?", invitation.WorkspaceID, userID).First(&existingMember).Error; err == nil {
		return echo.NewHTTPError(http.StatusConflict, "You are already a member of this workspace")
	}

	// Create workspace membership
	member := models.WorkspaceMember{
		WorkspaceID: invitation.WorkspaceID,
		UserID:      userID,
		Role:        invitation.Role,
	}

	if err := h.db.Create(&member).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add user to workspace")
	}

	// Mark invitation as accepted
	now := time.Now()
	invitation.AcceptedAt = &now
	if err := h.db.Save(&invitation).Error; err != nil {
		// Log error but don't fail the request since membership was created
		fmt.Printf("Warning: Failed to update invitation status: %v\n", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Successfully joined workspace",
		"workspace": invitation.Workspace,
		"role":      member.Role,
	})
}

// GetInvitations lists pending invitations for a workspace
func (h *WorkspaceHandler) GetInvitations(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workspace ID")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	// Check if user has permission to view invitations (owner or admin)
	var membership models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&membership).Error; err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	if membership.Role != "owner" && membership.Role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to view invitations")
	}

	// Get pending invitations
	var invitations []models.WorkspaceInvitation
	if err := h.db.Preload("Inviter").Where("workspace_id = ? AND accepted_at IS NULL AND expires_at > ?", workspaceID, time.Now()).Find(&invitations).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get invitations")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"invitations": invitations,
	})
}

// CancelInvitation cancels a pending invitation
func (h *WorkspaceHandler) CancelInvitation(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workspace ID")
	}

	invitationID, err := uuid.Parse(c.Param("invitation_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid invitation ID")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	// Check if user has permission to cancel invitations (owner or admin)
	var membership models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&membership).Error; err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	if membership.Role != "owner" && membership.Role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to cancel invitations")
	}

	// Find and delete the invitation
	var invitation models.WorkspaceInvitation
	if err := h.db.Where("id = ? AND workspace_id = ? AND accepted_at IS NULL", invitationID, workspaceID).First(&invitation).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invitation not found")
	}

	if err := h.db.Delete(&invitation).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to cancel invitation")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Invitation cancelled successfully",
	})
}

// GetMembers lists all members of a workspace
func (h *WorkspaceHandler) GetMembers(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workspace ID")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	// Check if user is a member of the workspace
	var membership models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&membership).Error; err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Get all members
	var members []models.WorkspaceMember
	if err := h.db.Preload("User").Where("workspace_id = ?", workspaceID).Find(&members).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get members")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"members": members,
	})
}

// UpdateMemberRoleRequest represents member role update request
type UpdateMemberRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin member viewer"`
}

// UpdateMemberRole updates a member's role in the workspace
func (h *WorkspaceHandler) UpdateMemberRole(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workspace ID")
	}

	memberID, err := uuid.Parse(c.Param("member_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid member ID")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	var req UpdateMemberRoleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check if user has permission to update roles (owner or admin)
	var membership models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&membership).Error; err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	if membership.Role != "owner" && membership.Role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to update member roles")
	}

	// Get the member to update
	var targetMember models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id = ?", workspaceID, memberID).First(&targetMember).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Member not found")
	}

	// Prevent changing owner role
	if targetMember.Role == "owner" {
		return echo.NewHTTPError(http.StatusForbidden, "Cannot change owner role")
	}

	// Prevent non-owners from promoting to admin or changing admin roles
	if membership.Role != "owner" && (req.Role == "admin" || targetMember.Role == "admin") {
		return echo.NewHTTPError(http.StatusForbidden, "Only owners can manage admin roles")
	}

	// Update role
	targetMember.Role = req.Role
	if err := h.db.Save(&targetMember).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update member role")
	}

	// Load updated member with user info
	if err := h.db.Preload("User").First(&targetMember, targetMember.ID).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load updated member")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"member":  targetMember,
		"message": "Member role updated successfully",
	})
}

// RemoveMember removes a member from the workspace
func (h *WorkspaceHandler) RemoveMember(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workspace ID")
	}

	memberID, err := uuid.Parse(c.Param("member_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid member ID")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	// Check if user has permission to remove members (owner or admin)
	var membership models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&membership).Error; err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	if membership.Role != "owner" && membership.Role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to remove members")
	}

	// Get the member to remove
	var targetMember models.WorkspaceMember
	if err := h.db.Where("workspace_id = ? AND user_id = ?", workspaceID, memberID).First(&targetMember).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Member not found")
	}

	// Prevent removing owner
	if targetMember.Role == "owner" {
		return echo.NewHTTPError(http.StatusForbidden, "Cannot remove workspace owner")
	}

	// Prevent non-owners from removing admins
	if membership.Role != "owner" && targetMember.Role == "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "Only owners can remove admins")
	}

	// Allow users to remove themselves (leave workspace)
	if userID == memberID {
		// Additional check: prevent owner from leaving
		if targetMember.Role == "owner" {
			return echo.NewHTTPError(http.StatusForbidden, "Owner cannot leave workspace. Transfer ownership first.")
		}
	}

	// Remove member
	if err := h.db.Delete(&targetMember).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove member")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Member removed successfully",
	})
}