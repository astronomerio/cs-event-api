package v1

import (
	"github.com/astronomerio/event-api/ingestion"
	"github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

// RouteHandler contains all event API endpoints
type RouteHandler struct {
	ingestionHandler ingestion.MessageWriter
}

// NewRouteHandler returns a new RouteHandler
func NewRouteHandler(writer ingestion.MessageWriter) *RouteHandler {
	return &RouteHandler{ingestionHandler: writer}
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

	v1Pixel := router.Group("v1/pixel").Use(limits.RequestSizeLimiter(15000))
	{
		v1Pixel.GET("track", h.pixelHandler("track"))
		v1Pixel.GET("page", h.pixelHandler("page"))
		v1Pixel.GET("alias", h.pixelHandler("alias"))
		v1Pixel.GET("identify", h.pixelHandler("identify"))
		v1Pixel.GET("group", h.pixelHandler("group"))
	}
}
