package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/cpu"
	"runtime"
	"time"
)

var (
	MemoryAlloc = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_mem_alloc_bytes",
			Help: "Количество байт которое выделилось и используется",
		})

	MemoryTotalAlloc = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_mem_total_alloc_bytes",
			Help: "Общее количество байт которое было выделено",
		})

	MemorySys = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_mem_sys_bytes",
			Help: "Общее количество байт которое было выделено системой",
		})

	NumGoroutine = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_goroutines_count",
			Help: "Количество горутин которые сейчас существуют",
		})

	GCPauseTotalNs = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_gc_pause_ns_total",
			Help: "Общее количество времени, которое потрачено на сборку мусора",
		})

	GCCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_gc_count",
			Help: "Количество сборок мусора",
		})

	GCNextTarget = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_gc_next_target_bytes",
			Help: "Целевое количество байт для следующей сборки мусора",
		})

	CPUUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "Использование CPU в процентах",
		})
)

func StartRuntimeMetricsCollector(interval time.Duration) {
	go func() {
		var memStats runtime.MemStats

		for {
			runtime.ReadMemStats(&memStats)

			// Память
			MemoryAlloc.Set(float64(memStats.Alloc))
			MemoryTotalAlloc.Set(float64(memStats.TotalAlloc))
			MemorySys.Set(float64(memStats.Sys))

			// Горутины
			NumGoroutine.Set(float64(runtime.NumGoroutine()))

			// GC
			GCCount.Set(float64(memStats.NumGC))
			GCPauseTotalNs.Set(float64(memStats.PauseTotalNs))
			GCNextTarget.Set(float64(memStats.NextGC))

			time.Sleep(interval)
		}
	}()
}

func StartCPUMetricsCollector(interval time.Duration) {
	go func() {
		for {
			percents, err := cpu.Percent(0, false)
			if err == nil && len(percents) > 0 {
				CPUUsage.Set(percents[0])
			}

			time.Sleep(interval)
		}
	}()
}
