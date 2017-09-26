package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/api/prometheus"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/api/routes"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/api/v1"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/logger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type APIServer struct {
	RouteHandlers []routes.RouteHandler

	router     *gin.Engine
	httpserver *http.Server

	adminRouter     *gin.Engine
	adminHttpserver *http.Server

	config *APIServerConfig

	healthy                bool
	shouldStartAdminServer bool
}

type APIServerConfig struct {
	APIPort   string
	AdminPort string

	IngestionHandler ingestion.IngestionHandler

	GracefulShutdownDelay int

	Logger logger.Logger
}

func NewServer() *APIServer {
	s := APIServer{
		router:                 gin.New(),
		adminRouter:            gin.New(),
		healthy:                false,
		shouldStartAdminServer: false,
	}
	return &s
}

// WithConfig sets the servers config
func (s *APIServer) WithConfig(config *APIServerConfig) *APIServer {
	s.config = config
	return s
}

// WithDefaultRoutes adds the default routes we will always want
func (s *APIServer) WithDefaultRoutes() *APIServer {
	s.RouteHandlers = append(s.RouteHandlers, v1.NewRouteHandler())
	return s
}

// WithHealthCheck creates a http route to report the health of the http server.
// Generally used to report a bad status when shutting down; to allow LB's to gracefully
// remove it from the pool
func (s *APIServer) WithHealthCheck() *APIServer {
	s.adminRouter.GET("/health", s.HealthCheckHandler)
	s.shouldStartAdminServer = true
	return s
}

// WithPProf injects a middleware handler for pprof on the admin router
func (s *APIServer) WithPProf() *APIServer {
	pprof.Register(s.adminRouter, nil)
	s.shouldStartAdminServer = true
	return s
}

// WithPrometheusMonitoring injects a middleware handler that will hook into the prometheus client
func (s *APIServer) WithPrometheusMonitoring() *APIServer {
	prometheus.Register(s.adminRouter, s.router)
	s.shouldStartAdminServer = true
	return s
}

func (s *APIServer) WithRequestID() *APIServer {
	s.router.Use(RequestIdMiddleware())
	return s
}

func buildListeningInterface(port string) string {
	return "0.0.0.0:" + port
}

// Run starts the http server(s) and then listens for the shutdown signal
func (s *APIServer) Run() {
	if os.ExpandEnv("GIN_MODE") == gin.ReleaseMode {
		gin.DisableConsoleColor()
	}

	s.httpserver = &http.Server{
		Addr:    buildListeningInterface(s.config.APIPort),
		Handler: s.router,
	}

	handlerConfig := &routes.HandlerConfig{
		IngestionHandler: s.config.IngestionHandler,
	}

	for _, handler := range s.RouteHandlers {
		handler.Register(s.router)
		handler.WithConfig(handlerConfig)
	}

	if s.shouldStartAdminServer {
		s.adminHttpserver = &http.Server{
			Addr:    buildListeningInterface(s.config.AdminPort),
			Handler: s.adminRouter,
		}
		go func() {
			if err := s.adminHttpserver.ListenAndServe(); err != nil {
				log.Fatalf("listen adminHttpserver: %s\n", err)
			}
		}()
	}

	go func() {
		if err := s.httpserver.ListenAndServe(); err != nil {
			log.Fatalf("listen httpserver: %s\n", err)
		}
	}()

	s.SetHealthy()

	quit := setupSignalHandler()

	<-quit

	log.Println("Shutdown signal received. Gracefully shutting down...")

	s.SetUnhealthy()
	sleepDuration := time.Duration(s.config.GracefulShutdownDelay) * time.Second

	log.Printf("Sleeping for %s...\n", sleepDuration.String())
	time.Sleep(sleepDuration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpserver.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	if s.shouldStartAdminServer {
		if err := s.adminHttpserver.Shutdown(ctx); err != nil {
			log.Fatal("admin http Server Shutdown:", err)
		}
	}
}
