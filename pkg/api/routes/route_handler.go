package routes

import (
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

type RouteHandler interface {
	WithConfig(*HandlerConfig)
	Register(*gin.Engine)
}

type HandlerConfig struct {
	IngestionHandler ingestion.IngestionHandler
	Logger           logger.Logger
}

type RouteDefintion struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}
