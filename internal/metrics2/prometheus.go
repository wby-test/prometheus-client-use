// todo: prometheus-client 提供的指标
package metrics2

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// serviceName 需要服务 deployment 含有 app 标签，或者通过 prometheus 标签操作，无需在此定义
	serviceName  = os.Getenv("app")
	commonLabels = []string{"path", "code", "host", "method", "service_code"}

	uptime = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "uptime",
			Help: "HTTP service uptime.",
		}, []string{"serviceName"},
	)

	http_request_total = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "The total number of processed http requests",
		},
		commonLabels,
	)

	http_request_duration_milliseconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_millisecond",
			Help:    "Histogram of lantencies for HTTP requests",
			Buckets: []float64{100, 200, 300, 400, 500, 600, 1000, 3000, 8000},
		},
		commonLabels,
	)

	http_request_in_flight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_request_in_flight",
			Help: "Current number of http requests in flight",
		},
		commonLabels,
	)

	http_request_summary_seconds = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_summary_seconds",
			Help: "Summary of lantencies for HTTP requests",
		},
		commonLabels,
	)
)

//func Init() {
//	prometheus.MustRegister(
//		uptime,
//		http_request_total,
//		http_request_duration_milliseconds,
//		http_request_in_flight,
//		http_request_summary_seconds,
//	)
//
//	go recordUptime()
//}

func recordUptime() {
	for range time.Tick(time.Second) {
		uptime.WithLabelValues(serviceName).Inc()
	}
}

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == metricsPath {
			c.Next()
			return
		}
		now := time.Now()
		c.Next()
		code := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		host := c.RemoteIP()
		// kuplus 服务使用
		serviceCode := c.GetString("code")
		labels := []string{path, code, host, method, serviceCode}
		http_request_total.WithLabelValues(labels...).Inc()
		http_request_duration_milliseconds.WithLabelValues(labels...).Observe(float64(time.Since(now).Milliseconds()))
		http_request_in_flight.WithLabelValues(labels...).Inc()
		defer http_request_in_flight.WithLabelValues(labels...).Dec()
		http_request_summary_seconds.WithLabelValues(labels...).Observe(time.Since(now).Seconds())
	}
}

const (
	metricsPath = "/metrics"
	healthzPath = "/healthz"
)
