package routes

import (
	"github.com/astronomerio/event-api/ingestion"
	"github.com/gin-gonic/gin"
)

// RouteHandler defines a generic type that can register gin routes
type RouteHandler interface {
	Register(*gin.Engine)
}

// RouteHandlerConfig defines configurations that can be applied to a RouteHandler
type RouteHandlerConfig struct {
	MessageWriter ingestion.MessageWriter
}
