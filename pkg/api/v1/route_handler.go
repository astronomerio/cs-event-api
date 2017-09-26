package v1

import (
	"github.com/astronomerio/clickstream-ingestion-api/pkg/api/routes"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

type RouteHandler struct {
	ingestionHandler ingestion.IngestionHandler
	logger           logger.Logger
}

func NewRouteHandler() *RouteHandler {
	return &RouteHandler{}
}

func (h *RouteHandler) WithConfig(config *routes.HandlerConfig) {
	h.ingestionHandler = config.IngestionHandler
	h.logger = config.Logger
}

func (h *RouteHandler) Register(router *gin.Engine) {
	v1Group := router.Group("v1")
	{
		v1Group.POST("t", h.singleHandler("track"))
		v1Group.GET("track", h.singleHandler("track"))

		v1Group.GET("p", h.singleHandler("page"))
		v1Group.GET("page", h.singleHandler("page"))

		v1Group.GET("a", h.singleHandler("alias"))
		v1Group.GET("alias", h.singleHandler("alias"))

		v1Group.GET("i", h.singleHandler("identify"))
		v1Group.GET("identify", h.singleHandler("identify"))

		v1Group.GET("g", h.singleHandler("group"))
		v1Group.GET("group", h.singleHandler("group"))

		v1Group.GET("import", h.importHandler)
		v1Group.GET("batch", h.batchHandler)
	}
}
