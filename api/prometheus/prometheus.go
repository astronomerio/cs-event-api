package prometheus

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "event_api"
)

var (
	// Create the request counter
	requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "event_api_requests_total",
		Help: "How many api requests processed, partitioned by type and write_key",
	}, []string{"method"})

	// Create the event counter
	eventCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "event_api_events_total",
		Help: "How many events processed, partitioned by write_key",
	}, []string{"method"})

	// Create the error counter
	errorCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "event_api_errors_total",
		Help: "How many errors from the API",
	}, []string{"method", "action"})

	// Create the duration histogram
	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "event_api_request_duration_seconds",
		Help:    "The API Request latencies in seconds",
		Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	}, []string{"method"})
)

func init() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(eventCounter)
	prometheus.MustRegister(errorCounter)
	prometheus.MustRegister(requestDuration)
}

// Register applies the route for prometheus scraping and applies the middleware function
// for profiling
func Register(router, middlewareRouter *gin.Engine) {
	// Register the /metrics route
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// Set up our middleware
	middlewareRouter.Use(middleware)
}

func middleware(ctx *gin.Context) {
	// Grab the current time
	start := time.Now()

	// Run the rest of the request
	ctx.Next()

	// Grab the values that were set in the handlers
	m := ctx.GetString("method")
	a := ctx.GetString("action")
	ec := ctx.GetInt("event_count")
	er := ctx.GetString("error")

	// Bail if no method
	if len(m) == 0 {
		return
	}

	// Increment errors if there was an error
	if len(er) > 0 {
		errorCounter.WithLabelValues(m, a).Inc()
	}

	// Update metrics
	requestCounter.WithLabelValues(m).Inc()
	eventCounter.WithLabelValues(m).Add(float64(ec))
	requestDuration.WithLabelValues(m).Observe(time.Since(start).Seconds())
}
