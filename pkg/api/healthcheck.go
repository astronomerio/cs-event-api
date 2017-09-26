package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler writes to the request whether or not the server is healthy
func (s *APIServer) HealthCheckHandler(c *gin.Context) {
	if s.IsHealthly() {
		c.AbortWithStatus(http.StatusOK)
		return
	}
	c.AbortWithStatus(http.StatusServiceUnavailable)
}

// SetHealthy marks the server as healthy
func (s *APIServer) SetHealthy() {
	s.healthy = true
	s.httpserver.SetKeepAlivesEnabled(true)
}

// SetUnhealthy marks the server as unhealthy
func (s *APIServer) SetUnhealthy() {
	s.healthy = false
	s.httpserver.SetKeepAlivesEnabled(false)
}

// IsHealthly returns whether or not the server is healthy
func (s *APIServer) IsHealthly() bool {
	return s.healthy
}
