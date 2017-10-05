package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler writes to the request whether or not the server is healthy
func (s *APIServer) HealthCheckHandler(c *gin.Context) {
	if s.IsHealthy() {
		c.AbortWithStatus(http.StatusOK)
	} else {
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// SetHealthy marks the server as healthy
func (s *APIServer) SetHealthy() {
	s.healthy = true
	s.httpServer.SetKeepAlivesEnabled(true)
}

// SetUnhealthy marks the server as unhealthy
func (s *APIServer) SetUnhealthy() {
	s.healthy = false
	s.httpServer.SetKeepAlivesEnabled(false)
}

// IsHealthly returns whether or not the server is healthy
func (s *APIServer) IsHealthy() bool {
	return s.healthy
}
