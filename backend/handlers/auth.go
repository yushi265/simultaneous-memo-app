package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
	
	"simultaneous-memo-app/backend/auth"
	"simultaneous-memo-app/backend/models"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token     string           `json:"token"`
	User      models.User      `json:"user"`
	Workspace models.Workspace `json:"workspace"`
}

// Register handles user registration
func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check if user already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", strings.ToLower(req.Email)).First(&existingUser).Error; err == nil {
		return echo.NewHTTPError(http.StatusConflict, "user already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
	}

	// Create user (personal workspace will be created automatically by the AfterCreate hook)
	user := models.User{
		Email:        strings.ToLower(req.Email),
		PasswordHash: hashedPassword,
		Name:         req.Name,
	}

	if err := h.db.Create(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}

	// Get the personal workspace
	var workspace models.Workspace
	if err := h.db.Where("owner_id = ? AND is_personal = ?", user.ID, true).First(&workspace).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get workspace")
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Name, workspace.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		Token:     token,
		User:      user,
		Workspace: workspace,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Find user
	var user models.User
	if err := h.db.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// Get user's personal workspace
	var workspace models.Workspace
	if err := h.db.Where("owner_id = ? AND is_personal = ?", user.ID, true).First(&workspace).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get workspace")
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Name, workspace.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token:     token,
		User:      user,
		Workspace: workspace,
	})
}

// Me returns current user information
func (h *AuthHandler) Me(c echo.Context) error {
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	workspaceID, ok := c.Get("workspaceID").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "workspace not found")
	}

	// Get user
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// Get current workspace
	var workspace models.Workspace
	if err := h.db.First(&workspace, workspaceID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "workspace not found")
	}

	// Get all user's workspaces
	var memberships []models.WorkspaceMember
	if err := h.db.Preload("Workspace").Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get workspaces")
	}

	workspaces := make([]map[string]interface{}, len(memberships))
	for i, membership := range memberships {
		workspaces[i] = map[string]interface{}{
			"id":   membership.Workspace.ID,
			"name": membership.Workspace.Name,
			"role": membership.Role,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":              user,
		"currentWorkspace":  workspace,
		"workspaces":        workspaces,
	})
}

// Logout handles user logout (client-side will remove the token)
func (h *AuthHandler) Logout(c echo.Context) error {
	// In a JWT-based system, logout is typically handled client-side
	// by removing the token. Here we just return a success response.
	return c.JSON(http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}