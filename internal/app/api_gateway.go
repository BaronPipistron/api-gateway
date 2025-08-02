package app

import (
	"context"
	"fmt"
	"github.com/BaronPipistron/api-gateway/internal/app/utils"
	"github.com/BaronPipistron/api-gateway/internal/config"
	"github.com/BaronPipistron/api-gateway/internal/telemetry/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	internalHttp "github.com/BaronPipistron/api-gateway/internal/presentation/http"
)

func Run() {
	cfg := loadConfig()

	ctx, stop := newSignalContext()
	defer stop()

	router := setupRouter(cfg)
	srv := createServer(cfg, router)
	startServer(srv)

	<-ctx.Done()
	shutdownServer(srv)
}

func loadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		panic(err)
	}
	logging.Init(cfg)
	return cfg
}

func newSignalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

func setupRouter(cfg *config.Config) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(utils.LoggingMiddleware(logging.Logger))

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	internalHttp.RegisterRoutes(engine, cfg)

	return engine
}

func createServer(cfg *config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.Server.HttpPort,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func startServer(srv *http.Server) {
	go func() {
		logging.Info("Starting server on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Fatal("Server failed:", err)
		}
	}()
}

func shutdownServer(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logging.Info("Graceful shutdown initialized...")

	if err := srv.Shutdown(ctx); err != nil {
		logging.Fatal("Error while server shutdown:", err)
	} else {
		logging.Info("Server stopped gracefully")
	}
}
