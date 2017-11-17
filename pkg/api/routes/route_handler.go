package routes

import (
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type RouteHandler interface {
	WithConfig(*HandlerConfig)
	Register(*gin.Engine)
}

type HandlerConfig struct {
	IngestionHandler ingestion.Handler
	Logger           *logrus.Logger
}

type RouteDefinition struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}
