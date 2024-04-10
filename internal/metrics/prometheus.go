package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	http_request_total = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "The total number of processed http requests",
		},
		[]string{"path", "code", "host", "method", "service_code"},
	)

	http_request_duration_milliseconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_millisecond",
			Help:    "Histogram of lantencies for HTTP requests",
			Buckets: []float64{100, 200, 300, 400, 500, 600, 1000, 3000, 8000},
		},
		[]string{"path", "code"},
	)

	//http_request_in_flight = promauto.NewGauge(
	//	prometheus.GaugeOpts{
	//		Name: "http_request_in_flight",
	//		Help: "Current number of http requests in flight",
	//	},
	//)
	//)

	//http_request_summary_seconds = promauto.NewSummary(
	//	prometheus.SummaryOpts{
	//		Name: "http_request_summary_seconds",
	//		Help: "Summary of lantencies for HTTP requests",
	//	},
	//)
)

const (
	metricsPath = "/metrics"
	healthzPath = "/healthz"
)

func AddHandleFunc(router *gin.Engine) {
	router.GET(healthzPath, func(c *gin.Context) { c.Status(http.StatusOK) })
	router.GET(metricsPath, prometheusHandler())
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func CountAndDuration() gin.HandlerFunc {
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
		serviceCode := c.GetString("code")
		http_request_total.WithLabelValues(path, code, host, method, serviceCode).Inc()
		http_request_duration_milliseconds.WithLabelValues(path, code).Observe(float64(time.Since(now).Milliseconds()))
		//http_request_in_flight.Inc()
		//defer http_request_in_flight.Dec()
		//http_request_summary_seconds.Observe(time.Since(now).Seconds())
	}
}
