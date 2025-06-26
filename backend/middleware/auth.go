package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	"simultaneous-memo-app/backend/auth"
)

// AuthMiddleware validates JWT tokens and adds user info to context
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// Check if the header starts with "Bearer "
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			tokenString := tokenParts[1]

			// Validate token
			claims, err := auth.ValidateToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Add user info to context
			c.Set("userID", claims.UserID)
			c.Set("userEmail", claims.Email)
			c.Set("userName", claims.Name)
			c.Set("workspaceID", claims.CurrentWorkspaceID)

			return next(c)
		}
	}
}

// GetUserID gets user ID from context
func GetUserID(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return uuid.UUID{}, echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}
	return userID, nil
}

// GetWorkspaceID gets workspace ID from context
func GetWorkspaceID(c echo.Context) (uuid.UUID, error) {
	workspaceID, ok := c.Get("workspaceID").(uuid.UUID)
	if !ok {
		return uuid.UUID{}, echo.NewHTTPError(http.StatusUnauthorized, "workspace not found")
	}
	return workspaceID, nil
}