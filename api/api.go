package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/astronomerio/event-api/api/prometheus"
	"github.com/astronomerio/event-api/api/routes"
	"github.com/astronomerio/event-api/api/v1"
	"github.com/astronomerio/event-api/ingestion"

	"os/signal"
	"syscall"

	"github.com/astronomerio/event-api/logging"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Server represents this server
type Server struct {
	RouteHandlers          []routes.RouteHandler
	router                 *gin.Engine
	httpServer             *http.Server
	adminRouter            *gin.Engine
	adminHTTPServer        *http.Server
	config                 *ServerConfig
	healthy                bool
	shouldStartAdminServer bool
}

// ServerConfig holds configurations for this server
type ServerConfig struct {
	APIPort               string
	AdminPort             string
	APIInterface          string
	AdminInterface        string
	MessageWriter         ingestion.MessageWriter
	GracefulShutdownDelay int
}

// NewServer creates a new server
func NewServer() *Server {
	// Create new server
	s := Server{
		router:                 gin.New(),
		adminRouter:            gin.New(),
		healthy:                false,
		shouldStartAdminServer: false,
	}

	// Set up middleware
	s.router.Use(gin.Recovery())
	return &s
}

// WithConfig sets the servers config
func (s *Server) WithConfig(config *ServerConfig) *Server {
	s.config = config
	return s
}

// WithDefaultRoutes adds the default routes we will always want
func (s *Server) WithDefaultRoutes() *Server {
	s.RouteHandlers = append(s.RouteHandlers, v1.NewRouteHandler())
	return s
}

// WithHealthCheck creates a http route to report the health of the http server.
// Generally used to report a bad status when shutting down; to allow LB's to gracefully
// remove it from the pool
func (s *Server) WithHealthCheck() *Server {
	s.adminRouter.GET("/healthz", s.HealthCheckHandler)
	s.shouldStartAdminServer = true
	return s
}

// WithPProf injects a middleware handler for pprof on the admin router
func (s *Server) WithPProf() *Server {
	pprof.Register(s.adminRouter, nil)
	s.shouldStartAdminServer = true
	return s
}

// WithPrometheusMonitoring injects a middleware handler that will hook into the prometheus client
func (s *Server) WithPrometheusMonitoring() *Server {
	prometheus.Register(s.adminRouter, s.router)
	s.shouldStartAdminServer = true
	return s
}

// WithRequestID injects request id middleware
func (s *Server) WithRequestID() *Server {
	s.router.Use(RequestIDMiddleware())
	return s
}

// Run starts the http server(s) and then listens for the shutdown signal
func (s *Server) Run() {
	log := logging.GetLogger(logrus.Fields{"package": "api"})

	if os.ExpandEnv("GIN_MODE") == gin.ReleaseMode {
		gin.DisableConsoleColor()
	}

	s.httpServer = &http.Server{
		Addr:    s.config.APIInterface + ":" + s.config.APIPort,
		Handler: s.router,
	}

	handlerConfig := &routes.RouteHandlerConfig{
		MessageWriter: s.config.MessageWriter,
	}

	handlerConfig.MessageWriter.Start()

	for _, handler := range s.RouteHandlers {
		handler.Register(s.router)
		handler.WithConfig(handlerConfig)
	}

	// Start admin server
	if s.shouldStartAdminServer {
		log.Info("Starting administrative server")

		s.adminHTTPServer = &http.Server{
			Addr:    s.config.AdminInterface + ":" + s.config.AdminPort,
			Handler: s.adminRouter,
		}

		go func() {
			if err := s.adminHTTPServer.ListenAndServe(); err != nil {
				log.Fatalf("Listen adminHTTPServer: %s\n", err)
			}
		}()
	}

	// Start events server
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("listen httpserver: %s\n", err)
		}
	}()

	s.SetHealthy()

	c := make(chan os.Signal)
	signal.Notify(c,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGSTOP,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGKILL)

	log.Info(<-c)
	s.stop(handlerConfig.MessageWriter)
	os.Exit(1)
}

func (s *Server) stop(writer ingestion.MessageWriter) {
	log := logging.GetLogger(logrus.Fields{"package": "api"})
	log.Info("Shutdown signal received. Gracefully shutting down...")
	s.SetUnhealthy()
	sleepDuration := time.Duration(s.config.GracefulShutdownDelay) * time.Second

	log.Infof("Sleeping for %s...", sleepDuration.String())
	time.Sleep(sleepDuration)

	err := writer.Shutdown()
	if err != nil {
		log.Errorf("error shutting down ingestion handler %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Errorf("server shutdown: %s", err.Error())
	}

	if s.shouldStartAdminServer {
		if err := s.adminHTTPServer.Shutdown(ctx); err != nil {
			log.Errorf("admin http server shutdown: %s", err.Error())
		}
	}

}
