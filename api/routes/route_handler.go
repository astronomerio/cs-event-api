package routes

import (
	"github.com/astronomerio/event-api/ingestion"

	"github.com/gin-gonic/gin"
)

// RouteHandler defines a generic type that can register gin routes
type RouteHandler interface {
	// WithConfig(*RouteHandlerConfig)
	Register(*gin.Engine)
}

// RouteHandlerConfig defines configurations that can be applied to a RouteHandler
type RouteHandlerConfig struct {
	MessageWriter ingestion.MessageWriter
	// Logger        *logrus.Logger
}

// RouteDefinition defines a gin route definition
// type RouteDefinition struct {
// 	Method  string
// 	Path    string
// 	Handler gin.HandlerFunc
// }
