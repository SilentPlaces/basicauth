package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
	httpInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_in_flight_requests",
			Help: "Number of in-flight HTTP requests.",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, httpInFlight)
}

func PrometheusMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpInFlight.Inc()
		start := time.Now()

		c.Next()

		httpInFlight.Dec()
		statusCode := strconv.Itoa(c.Writer.Status())
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		httpRequestsTotal.WithLabelValues(c.Request.Method, route, statusCode).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, route, statusCode).Observe(time.Since(start).Seconds())
	}
}
