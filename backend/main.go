package main

import (
	"log"
	"net/http"
	"simultaneous-memo-app/backend/config"
	"simultaneous-memo-app/backend/handlers"
	"simultaneous-memo-app/backend/models"
	"simultaneous-memo-app/backend/websocket"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize handlers
	h := handlers.NewHandler(db)

	// Routes
	api := e.Group("/api")
	
	// Page routes
	api.GET("/pages", h.GetPages)
	api.POST("/pages", h.CreatePage)
	api.GET("/pages/:id", h.GetPage)
	api.PUT("/pages/:id", h.UpdatePage)
	api.DELETE("/pages/:id", h.DeletePage)

	// Image upload
	api.POST("/upload", h.UploadFile)
	
	// General file upload
	api.POST("/upload/file", h.UploadGeneralFile)
	api.GET("/files", h.ListFiles)
	api.GET("/files/:id", h.GetFileMetadata)
	api.DELETE("/files/:id", h.DeleteFile)
	api.GET("/files/*", h.GetFile)
	api.GET("/file/*", h.ServeFile)

	// Image management
	api.GET("/images", h.GetImages)
	api.GET("/images/:id", h.GetImageByID)
	api.DELETE("/images/:id", h.DeleteImageByID)
	
	// Responsive image serving
	api.GET("/img/*", h.ServeImage)
	
	// Admin endpoints
	api.POST("/admin/cleanup-images", h.CleanupImages)

	// WebSocket endpoint
	ws := websocket.NewHub()
	go ws.Run()
	
	e.GET("/ws/:pageId", func(c echo.Context) error {
		pageID := c.Param("pageId")
		websocket.HandleWebSocket(ws, c.Response(), c.Request(), pageID)
		return nil
	})

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