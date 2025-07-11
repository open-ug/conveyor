package metrics

import (
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v4/cpu"
)

// Prometheus metrics
var (
	// HTTP metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)

	// System metrics
	memoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
		[]string{"type"},
	)

	cpuUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "CPU usage percentage",
		},
	)

	goroutinesCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines_count",
			Help: "Number of goroutines",
		},
	)

	// Application metrics
	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)
)

func init() {
	// Register all metrics
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		httpRequestsInFlight,
		memoryUsage,
		cpuUsage,
		goroutinesCount,
		activeConnections,
	)
}

// PrometheusMiddleware creates a middleware to collect HTTP metrics
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Increment in-flight requests
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		// Process request
		err := c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		method := c.Method()
		path := c.Route().Path
		if path == "" {
			path = "unknown"
		}
		statusCode := strconv.Itoa(c.Response().StatusCode())

		httpRequestsTotal.WithLabelValues(method, path, statusCode).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}

func getCurrentCPUUsage() float64 {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
		return 0.0
	}

	if len(percentages) > 0 {
		return percentages[0]
	}

	return 0.0
}

// updateSystemMetrics updates system-level metrics
func updateSystemMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Memory metrics
	memoryUsage.WithLabelValues("heap_alloc").Set(float64(m.HeapAlloc))
	memoryUsage.WithLabelValues("heap_sys").Set(float64(m.HeapSys))
	memoryUsage.WithLabelValues("heap_idle").Set(float64(m.HeapIdle))
	memoryUsage.WithLabelValues("heap_inuse").Set(float64(m.HeapInuse))
	memoryUsage.WithLabelValues("stack_sys").Set(float64(m.StackSys))

	// Goroutines
	goroutinesCount.Set(float64(runtime.NumGoroutine()))

	cpuUsage.Set(getCurrentCPUUsage())
}

// startMetricsUpdater starts a goroutine to periodically update metrics
func StartAPIMetricsUpdater() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			updateSystemMetrics()
		}
	}()
}
