package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/api/prometheus"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/api/routes"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/api/v1"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion"

	"os/signal"
	"syscall"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/logging"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	RouteHandlers []routes.RouteHandler

	router     *gin.Engine
	httpServer *http.Server

	adminRouter     *gin.Engine
	adminHttpServer *http.Server

	config *ServerConfig

	healthy                bool
	shouldStartAdminServer bool
}

type ServerConfig struct {
	APIPort   string
	AdminPort string

	APIInterface   string
	AdminInterface string

	IngestionHandler ingestion.Handler

	GracefulShutdownDelay int
}

func NewServer() *Server {
	s := Server{
		router:                 gin.New(),
		adminRouter:            gin.New(),
		healthy:                false,
		shouldStartAdminServer: false,
	}
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
	s.adminRouter.GET("/health", s.HealthCheckHandler)
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

func (s *Server) WithRequestID() *Server {
	s.router.Use(RequestIdMiddleware())
	return s
}

// Run starts the http server(s) and then listens for the shutdown signal
func (s *Server) Run() {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "api", "function": "Run"})

	if os.ExpandEnv("GIN_MODE") == gin.ReleaseMode {
		gin.DisableConsoleColor()
	}

	s.httpServer = &http.Server{
		Addr:    s.config.APIInterface + ":" + s.config.APIPort,
		Handler: s.router,
	}

	handlerConfig := &routes.HandlerConfig{
		IngestionHandler: s.config.IngestionHandler,
		Logger:           logging.GetLogger(),
	}

	handlerConfig.IngestionHandler.Start()

	for _, handler := range s.RouteHandlers {
		handler.Register(s.router)
		handler.WithConfig(handlerConfig)
	}

	if s.shouldStartAdminServer {
		logger.Info("starting administrative server")
		s.adminHttpServer = &http.Server{
			Addr:    s.config.AdminInterface + ":" + s.config.AdminPort,
			Handler: s.adminRouter,
		}
		go func() {
			if err := s.adminHttpServer.ListenAndServe(); err != nil {
				logger.Fatalf("listen adminHttpServer: %s\n", err)
			}
		}()
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			logger.Fatalf("listen httpserver: %s\n", err)
		}
	}()

	s.SetHealthy()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL)
	logger.Info(<-c)
	s.stop(handlerConfig.IngestionHandler)
	os.Exit(1)
}

func (s *Server) stop(ingestionHandler ingestion.Handler) {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "api", "function": "stop"})
	logger.Info("Shutdown signal received. Gracefully shutting down...")
	s.SetUnhealthy()
	sleepDuration := time.Duration(s.config.GracefulShutdownDelay) * time.Second

	logger.Infof("Sleeping for %s...", sleepDuration.String())
	time.Sleep(sleepDuration)

	err := ingestionHandler.Shutdown()
	if err != nil {
		logger.Errorf("error shutting down ingestion handler %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("server shutdown: %s", err.Error())
	}

	if s.shouldStartAdminServer {
		if err := s.adminHttpServer.Shutdown(ctx); err != nil {
			logger.Errorf("admin http server shutdown: %s", err.Error())
		}
	}

}
