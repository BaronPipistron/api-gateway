package metrics

import (
	"fmt"
	"github.com/BaronPipistron/api-gateway/internal/config"
	"github.com/BaronPipistron/api-gateway/internal/telemetry/logging"
	"github.com/BaronPipistron/api-gateway/internal/telemetry/logging "
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func RegisterPrometheusHandler(cfg *config.Config) {
	http.Handle("/metrics", promhttp.Handler())

	metricsPort := cfg.Metrics.Port
	addr := fmt.Sprintf(":%d", metricsPort)

	logging.Logger.Info("Prometheus handler registered", "path", "/metrics", "port", metricsPort)
	go func() {
		logging.Logger.Info("Cтарт сервера метрик Prometheus", "address", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			logging.Logger.Fatal("Ошибка старта сервера метрик Prometheus", "error", err)
		}
	}()
}
