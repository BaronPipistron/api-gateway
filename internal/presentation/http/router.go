package http

import (
	"github.com/BaronPipistron/api-gateway/internal/app/utils"
	"github.com/BaronPipistron/api-gateway/internal/config"
	"github.com/BaronPipistron/api-gateway/internal/presentation/http/proxy"
	"github.com/BaronPipistron/api-gateway/internal/telemetry/logging"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, cfg *config.Config) {
	logger := logging.Logger

	router.Use(utils.LoggingMiddleware(logger))

	handler := proxy.NewProxyHandler(cfg)

	router.Any("/api/*any", handler.Handle)
	// router.NoRoute(handler.Handle) - чтобы обрабатывать все пути
}
