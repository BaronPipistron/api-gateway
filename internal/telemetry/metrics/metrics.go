package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var HTTPRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Общее количество HTTP запросов, прошедших через API Gateway",
	},
	[]string{"method", "path", "status"},
)
var HTTPResponseTime = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_response_time",
		Help:    "Время ответа HTTP запросов, прошедших через API Gateway",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path", "status"},
)

func IncrementHTTPRequestsTotal(method, path, status string) {
	HTTPRequestsTotal.With(prometheus.Labels{
		"method": method,
		"path":   path,
		"status": status,
	}).Inc()
}

func ObserveHTTPResponseTime(method, path string, duration float64) {
	HTTPResponseTime.With(prometheus.Labels{
		"method": method,
		"path":   path,
	}).Observe(duration)
}
