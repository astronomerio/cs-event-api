package v1

import (
	"github.com/astronomerio/event-api/pkg/api/routes"
	"github.com/astronomerio/event-api/pkg/ingestion"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RouteHandler struct {
	ingestionHandler ingestion.Handler
	logger           *logrus.Logger
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
		v1Group.POST("track", h.singleHandler("track"))

		v1Group.POST("p", h.singleHandler("page"))
		v1Group.POST("page", h.singleHandler("page"))

		v1Group.POST("a", h.singleHandler("alias"))
		v1Group.POST("alias", h.singleHandler("alias"))

		v1Group.POST("i", h.singleHandler("identify"))
		v1Group.POST("identify", h.singleHandler("identify"))

		v1Group.POST("g", h.singleHandler("group"))
		v1Group.POST("group", h.singleHandler("group"))

		v1Group.POST("import", h.importHandler)
		v1Group.POST("batch", h.batchHandler)
	}
}
