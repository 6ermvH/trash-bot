package panel

import (
	"context"
	"fmt"

	"github.com/6ermvH/trash-bot/internal/config"
	"github.com/gin-gonic/gin"
)

func Start(ctx context.Context, cfg *config.Config) error {
	router := gin.Default()

	// TODO router.GET(...) {...}

	port := fmt.Sprintf(":%s", cfg.Server.Port)
	return router.Run(port)
}
