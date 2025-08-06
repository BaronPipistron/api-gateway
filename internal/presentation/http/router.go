package http

import (
	"github.com/BaronPipistron/api-gateway/internal/app/utils"
	"github.com/BaronPipistron/api-gateway/internal/config"
	"github.com/BaronPipistron/api-gateway/internal/presentation/http/middleware"
	"github.com/BaronPipistron/api-gateway/internal/presentation/http/proxy"
	"github.com/BaronPipistron/api-gateway/internal/telemetry/logging"
	"github.com/BaronPipistron/api-gateway/internal/telemetry/metrics"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func NewRouter(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(requestLogger)
	r.Use(metricsMiddleware)

	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(cfg)

	r.Route("/internal", func(r chi.Router) {
		r.Use(apiKeyMiddleware.ValidateAPIKey)
		r.Get("/data", InternalDataHandler)
	})

	r.Get("/health", HealthCheckHandler)

	return r
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)

		logging.Logger.Info("Request processed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", duration.Seconds(),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
	})
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(startTime)
		statusCode := wrappedWriter.status
		method := r.Method
		path := r.URL.Path

		metrics.IncrementHTTPRequestsTotal(method, path, string(statusCode))
		metrics.ObserveHTTPResponseTime(method, path, duration.Seconds())
	})
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func InternalDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Internal Data - требует API Key"))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
func RegisterRoutes(router *gin.Engine, cfg *config.Config) {
	logger := logging.Logger

	router.Use(utils.LoggingMiddleware(logger))

	handler := proxy.NewProxyHandler(cfg)

	router.Any("/api/*any", handler.Handle)
	// router.NoRoute(handler.Handle) - чтобы обрабатывать все пути
}
