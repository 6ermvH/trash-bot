package panel

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/6ermvH/trash-bot/internal/config"
	handlers "github.com/6ermvH/trash-bot/internal/handlers/http/v1"
	"github.com/6ermvH/trash-bot/internal/services/trashmanager"
	"github.com/gin-gonic/gin"
)

//go:embed all:web
var webFS embed.FS

func serveEmbeddedFile(router *gin.Engine, route, path, contentType string) {
	router.GET(route, func(c *gin.Context) {
		data, err := webFS.ReadFile(path)
		if err != nil {
			log.Printf("read embedded file %s: %v", path, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load static file"})
			return
		}
		c.Data(http.StatusOK, contentType, data)
	})
}

func Start(ctx context.Context, cfg *config.Config, trashm *trashmanager.Service) error {
	router := gin.Default()
	router.RedirectTrailingSlash = false

	// API routes
	api := router.Group("/api")

	authHandler := handlers.NewAuthHandler(
		cfg.Server.AdminLogin,
		cfg.Server.AdminPassword,
		cfg.Server.JWTSecret,
	)
	api.POST("/login", authHandler.Login)

	handle := handlers.New(trashm)

	protected := api.Group("/")
	protected.Use(handlers.AuthMiddleware(cfg.Server.JWTSecret))
	{
		protected.GET("/stats", handle.Stats)
		protected.GET("/chats", handle.Chats)
		protected.GET("/chats/:id", handle.ChatByID)
	}

	// Static files
	serveEmbeddedFile(router, "/", "web/index.html", "text/html; charset=utf-8")
	serveEmbeddedFile(router, "/style.css", "web/style.css", "text/css; charset=utf-8")
	serveEmbeddedFile(router, "/app.js", "web/app.js", "application/javascript; charset=utf-8")

	port := ":" + cfg.Server.Port

	if err := router.Run(port); err != nil {
		return fmt.Errorf("Start server on port %s : %w", port, err)
	}

	return nil
}
