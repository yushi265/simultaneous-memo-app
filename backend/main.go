package main

import (
	"log"
	"net/http"
	"simultaneous-memo-app/backend/config"
	"simultaneous-memo-app/backend/handlers"
	"simultaneous-memo-app/backend/models"
	"simultaneous-memo-app/backend/websocket"
	customMiddleware "simultaneous-memo-app/backend/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/go-playground/validator/v10"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := models.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	if err := models.AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Echo
	e := echo.New()
	
	// Set custom validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize handlers
	h := handlers.NewHandler(db)
	authHandler := handlers.NewAuthHandler(db)
	workspaceHandler := handlers.NewWorkspaceHandler(db)

	// Initialize rate limiters
	fileUploadLimiter := customMiddleware.FileUploadRateLimiter()
	generalAPILimiter := customMiddleware.GeneralAPIRateLimiter()

	// Routes
	api := e.Group("/api")
	
	// Apply general rate limiting to all API routes
	api.Use(generalAPILimiter.Middleware())
	
	// Auth routes (public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
	auth.GET("/me", authHandler.Me, customMiddleware.AuthMiddleware())
	
	// Protected routes group
	protected := api.Group("")
	protected.Use(customMiddleware.AuthMiddleware())
	
	// Workspace routes (protected)
	protected.GET("/workspaces", workspaceHandler.GetWorkspaces)
	protected.POST("/workspaces", workspaceHandler.CreateWorkspace)
	protected.GET("/workspaces/:id", workspaceHandler.GetWorkspace)
	protected.PUT("/workspaces/:id", workspaceHandler.UpdateWorkspace)
	protected.DELETE("/workspaces/:id", workspaceHandler.DeleteWorkspace)
	protected.POST("/workspaces/:id/switch", workspaceHandler.SwitchWorkspace)
	
	// Page routes (protected)
	protected.GET("/pages", h.GetPages)
	protected.POST("/pages", h.CreatePage)
	protected.GET("/pages/:id", h.GetPage)
	protected.PUT("/pages/:id", h.UpdatePage)
	protected.DELETE("/pages/:id", h.DeletePage)

	// Image upload with stricter rate limiting (protected)
	protected.POST("/upload", h.UploadFile, fileUploadLimiter.Middleware())
	
	// General file upload with stricter rate limiting (protected)
	protected.POST("/upload/file", h.UploadGeneralFile, fileUploadLimiter.Middleware())
	protected.GET("/files", h.ListFiles)
	protected.GET("/files/:id", h.GetFileMetadata)
	protected.DELETE("/files/:id", h.DeleteFile)
	protected.GET("/file/*", h.ServeFile)

	// Image management (protected)
	protected.GET("/images", h.GetImages)
	protected.GET("/images/:id", h.GetImageByID)
	protected.DELETE("/images/:id", h.DeleteImageByID)
	
	// Responsive image serving (protected)
	protected.GET("/img/*", h.ServeImage)
	
	// Admin endpoints (protected)
	protected.POST("/admin/cleanup-images", h.CleanupImages)

	// WebSocket endpoint (protected)
	ws := websocket.NewHub()
	go ws.Run()
	
	e.GET("/ws/:pageId", func(c echo.Context) error {
		pageID := c.Param("pageId")
		// TODO: Add auth validation for WebSocket
		websocket.HandleWebSocket(ws, c.Response(), c.Request(), pageID)
		return nil
	}, customMiddleware.AuthMiddleware())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// CustomValidator implements echo.Validator interface
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the input
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}