package v1

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"

	v1types "github.com/astronomerio/clickstream-ingestion-api/pkg/types/v1"
	"github.com/gin-gonic/gin"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/logging"
	"github.com/sirupsen/logrus"
)

var returnJSON = map[string]bool{
	"success": true,
}

func (h *RouteHandler) singleHandler(kind string) gin.HandlerFunc {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "v1", "function": "singleHandler"})
	return func(c *gin.Context) {
		c.Set("profile", true)
		c.Set("type", "single")
		c.Set("action", kind)

		logger.Infof("Headers : %s", c.Request.Header)
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
