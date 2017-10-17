package prometheus

import (
	"time"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type prometheusInstrumentation struct {
	requestCounter  *prometheus.CounterVec
	errorCounter    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

var pi *prometheusInstrumentation

// Register applies the route for prometheus scraping and applies the middleware function
// for profiling
func Register(router, middlewareRouter *gin.Engine) {
	buildVectors()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	middlewareRouter.Use(middleware)
}

func buildVectors() {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "prometheus", "function": "buildVectors"})

	pi = &prometheusInstrumentation{
		requestCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "How many api requests processed, partitioned by type and action",
		}, []string{"type", "action"}),
		errorCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "api_errors_total",
			Help: "How many errors from the API",
		}, []string{"type", "action"}),
		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "The API Request latencies in seconds",
			Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		}, []string{"type", "action"}),
	}

	var err error
	if err = prometheus.Register(pi.requestCounter); err != nil {
		logger.Fatal("Error registering requestCounter", "Error", err)
	}
	if err = prometheus.Register(pi.requestDuration); err != nil {
		logger.Fatal("Error registering requestDuration", "Error", err)
	}
	if err = prometheus.Register(pi.errorCounter); err != nil {
		logger.Fatal("Error registering errorCounter", "Error", err)
	}
}

func middleware(ctx *gin.Context) {
	start := time.Now()

	ctx.Next()

	// the handler didnt mark this request to be profiled
	if !ctx.GetBool("profile") {
		return
	}

	t := ctx.GetString("type")
	a := ctx.GetString("action")

	if t == "" || a == "" {
		return
	}
	elapsed := time.Since(start)
	pi.requestCounter.WithLabelValues(t, a).Inc()
	pi.requestDuration.WithLabelValues(t, a).Observe(elapsed.Seconds())
}
