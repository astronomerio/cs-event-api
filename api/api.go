package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/astronomerio/event-api/api/prometheus"
	"github.com/astronomerio/event-api/api/routes"
	"github.com/gin-contrib/pprof"

	"github.com/astronomerio/event-api/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Server represents this server
type Server struct {
	handlers    []routes.RouteHandler
	server      *http.Server
	router      *gin.Engine
	adminServer *http.Server
	adminRouter *gin.Engine
	config      *ServerConfig
	healthy     bool
}

// ServerConfig holds configurations for this server
type ServerConfig struct {
	APIPort               string
	AdminPort             string
	APIInterface          string
	AdminInterface        string
	GracefulShutdownDelay int
}

// NewServer creates a new server
func NewServer(config *ServerConfig) *Server {
	// Create new server
	s := &Server{
		healthy:     false,
		router:      gin.Default(),
		adminRouter: gin.Default(),
	}

	// Set the config
	s.config = config

	// Create the actual http server
	s.server = &http.Server{
		Addr:    config.APIInterface + ":" + config.APIPort,
		Handler: s.router,
	}

	// Create the admin http server
	s.adminServer = &http.Server{
		Addr:    s.config.AdminInterface + ":" + s.config.AdminPort,
		Handler: s.adminRouter,
	}

	return s
}

// WithRouteHandler appends a new RouteHandler
func (s *Server) WithRouteHandler(rh routes.RouteHandler) *Server {
	s.handlers = append(s.handlers, rh)
	return s
}

// WithPProf injects a middleware handler for pprof on the admin router
func (s *Server) WithPProf() *Server {
	pprof.Register(s.adminRouter, nil)
	return s
}

// WithPrometheusMonitoring injects a middleware handler that will hook into the prometheus client
func (s *Server) WithPrometheusMonitoring() *Server {
	prometheus.Register(s.adminRouter, s.router)
	return s
}

// Serve starts the http server(s) and then listens for the shutdown signal
func (s *Server) Serve(shutdownChan <-chan struct{}) {
	log := logging.GetLogger(logrus.Fields{"package": "api"})

	if os.ExpandEnv("GIN_MODE") == gin.ReleaseMode {
		gin.DisableConsoleColor()
	}

	s.router.Use(RequestIDMiddleware())
	for _, handler := range s.handlers {
		handler.Register(s.router)
	}

	// Start admin server
	go func() {
		s.adminRouter.GET("/healthz", s.HealthCheckHandler)
		if err := s.adminServer.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	// Start events server
	go func() {
		s.SetHealthy()
		if err := s.server.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	<-shutdownChan
	log.Info("Webserver recieved shutdown signal")
}

// Close cleans up and shuts down the webservers
func (s *Server) Close() {
	log := logging.GetLogger(logrus.Fields{"package": "api"})
	s.SetUnhealthy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("Event server shutdown: %s", err)
	}

	if err := s.adminServer.Shutdown(ctx); err != nil {
		log.Errorf("Admin server shutdown: %s", err)
	}

	log.Info("Webserver has been shut down")
}
