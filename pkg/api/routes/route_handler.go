package routes

import (
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RouteHandler interface {
	WithConfig(*HandlerConfig)
	Register(*gin.Engine)
}

type HandlerConfig struct {
	IngestionHandler ingestion.IngestionHandler
	Log *logrus.Logger
}

type RouteDefinition struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}
