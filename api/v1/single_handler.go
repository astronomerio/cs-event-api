package v1

import (
	"net/http"
	"time"

	v1types "github.com/astronomerio/event-api/types/v1"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *RouteHandler) singleHandler(kind string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.logger.WithFields(logrus.Fields{"package": "v1"})
		c.Set("method", "single")

		// Create a new msg
		raw, err := c.GetRawData()
		if err != nil {
			// Log and return 200
			c.Set("error", err.Error())
			log.WithFields(logrus.Fields{
				"action": "read-body",
			}).Error(err.Error())
			c.AbortWithStatusJSON(http.StatusOK, returnJSON)
			return
		}

		// Unmarshal data into a Batch
		msg, err := v1types.NewMessage(kind, raw)
		if err != nil {
			// Log and return 200
			c.Set("error", err.Error())
			log.WithFields(logrus.Fields{
				"action": "single-unmarshal ",
			}).Error(err.Error())
			c.AbortWithStatusJSON(http.StatusOK, returnJSON)
			return
		}

		// Grab metadata from this request
		metadata := v1types.NewRequestMetadata(c)

		// Apply ReceivedAt date
		msg.WithReceivedAt(time.Now().UTC())

		// Apply metadata from context
		msg.WithRequestMetadata(metadata)

		// Skew timestamp to account for bad client clocks
		msg.SkewTimestamp()

		// Pass the msg along to the adapter
		h.ingestionHandler.ProcessMessage(msg.String(), msg.GetMessageID())

		// Set additional metric data
		c.Set("event_count", 1)

		// Always return 200
		c.Header("Connection", "keep-alive")
		c.AbortWithStatusJSON(http.StatusOK, returnJSON)
	}
}
