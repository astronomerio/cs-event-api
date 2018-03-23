package v1

import (
	"github.com/astronomerio/event-api/api/routes"
	"github.com/astronomerio/event-api/ingestion"
	"github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RouteHandler contains all event API endpoints
type RouteHandler struct {
	ingestionHandler ingestion.MessageWriter
	logger           *logrus.Logger
}

// NewRouteHandler returns a new RouteHandler
func NewRouteHandler() *RouteHandler {
	return &RouteHandler{}
}

// WithConfig applies a given config
func (h *RouteHandler) WithConfig(config *routes.RouteHandlerConfig) {
	h.ingestionHandler = config.MessageWriter
	h.logger = config.Logger
}

// Register registers the event handlers on the given router
func (h *RouteHandler) Register(router *gin.Engine) {
	v1Single := router.Group("v1").Use(limits.RequestSizeLimiter(15000))
	{
		v1Single.POST("t", h.singleHandler("track"))
		v1Single.POST("track", h.singleHandler("track"))

		v1Single.POST("p", h.singleHandler("page"))
		v1Single.POST("page", h.singleHandler("page"))

		v1Single.POST("a", h.singleHandler("alias"))
		v1Single.POST("alias", h.singleHandler("alias"))

		v1Single.POST("i", h.singleHandler("identify"))
		v1Single.POST("identify", h.singleHandler("identify"))

		v1Single.POST("g", h.singleHandler("group"))
		v1Single.POST("group", h.singleHandler("group"))
	}

	v1Batch := router.Group("v1").Use(limits.RequestSizeLimiter(500000))
	{
		v1Batch.POST("batch", h.batchHandler)
		v1Batch.POST("import", h.batchHandler)
	}
}
