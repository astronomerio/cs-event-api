package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *RouteHandler) importHandler(c *gin.Context) {
	c.Set("type", "batch")
	c.Set("action", "import")
	h.ingestionHandler.ProcessMessage("NOT IMPLEMENTED", "NOT IMPLEMENTED")
	c.AbortWithStatus(http.StatusOK)
}
