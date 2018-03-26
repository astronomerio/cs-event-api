package prometheus

import (
	"fmt"
	"time"

	"github.com/astronomerio/event-api/logging"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	requestCounter  *prometheus.CounterVec
	eventCounter    *prometheus.CounterVec
	errorCounter    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
)

const (
	namespace = "event_api"
)

func metricName(metric string) string {
	return fmt.Sprintf("%s_%s", namespace, metric)
}

// Register applies the route for prometheus scraping and applies the middleware function
// for profiling
func Register(router, middlewareRouter *gin.Engine) {
	// Build our prometheus vectors
	buildVectors()
	// Register the /metrics route
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// Set up our middleware
	middlewareRouter.Use(middleware)
}

func buildVectors() {
	var err error

	log := logging.GetLogger().WithFields(logrus.Fields{"package": "prometheus"})

	// Create the request counter
	requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metricName("requests_total"),
		Help: "How many api requests processed, partitioned by type and write_key",
	}, []string{"method"})

	if err = prometheus.Register(requestCounter); err != nil {
		log.Fatal("Error registering requestCounter:", err)
	}

	// Create the event counter
	eventCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metricName("events_total"),
		Help: "How many events processed, partitioned by write_key",
	}, []string{"method"})

	if err = prometheus.Register(eventCounter); err != nil {
		log.Fatal("Error registering eventCounter:", err)
	}

	// Create the error counter
	errorCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metricName("errors_total"),
		Help: "How many errors from the API",
	}, []string{"method", "action"})

	if err = prometheus.Register(errorCounter); err != nil {
		log.Fatal("Error registering errorCounter:", err)
	}

	// Create the duration histogram
	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    metricName("request_duration_seconds"),
		Help:    "The API Request latencies in seconds",
		Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	}, []string{"method"})

	if err = prometheus.Register(requestDuration); err != nil {
		log.Fatal("Error registering requestDuration:", err)
	}

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
