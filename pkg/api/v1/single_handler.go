package v1

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"

	v1types "github.com/astronomerio/cs-event-api/pkg/types/v1"
	"github.com/gin-gonic/gin"
)

var returnJSON = map[string]bool{
	"success": true,
}

func (h *RouteHandler) singleHandler(kind string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("profile", true)
		c.Set("type", "single")
		c.Set("action", kind)

		var message v1types.Message
		if err := c.ShouldBindWith(&message, binding.JSON); err != nil {
			c.Set("error", err.Error())
			c.Set("stage", "1")
			c.AbortWithStatusJSON(http.StatusOK, returnJSON)
			return
		}

		message.BindRequest(c)

		if !message.IsValid() {
			c.AbortWithStatusJSON(http.StatusOK, returnJSON)
			return
		}

		message.FormatTimestamps()
		message.MaybeFix()
		message.SkewTimestamp()

		h.ingestionHandler.ProcessMessage(message.String(), message.PartitionKey())

		c.Header("Connection", "keep-alive")
		c.AbortWithStatusJSON(http.StatusOK, returnJSON)
	}
}
