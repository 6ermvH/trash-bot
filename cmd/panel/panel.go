package panel

import (
	"context"
	"fmt"

	"github.com/6ermvH/trash-bot/internal/config"
	handlers "github.com/6ermvH/trash-bot/internal/handlers/rest/v1"
	"github.com/gin-gonic/gin"
)

func Start(ctx context.Context, cfg *config.Config) error {
	router := gin.Default()

	// TODO router.GET(...) {...}
	router.GET("/", handlers.Stat)

	port := ":" + cfg.Server.Port

	if err := router.Run(port); err != nil {
		return fmt.Errorf("Start server on port %s : %w", port, err)
	}

	return nil
}
