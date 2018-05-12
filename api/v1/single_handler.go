package v1

import (
	"net/http"
	"time"

	"github.com/astronomerio/event-api/logging"
	v1types "github.com/astronomerio/event-api/types/v1"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *RouteHandler) singleHandler(kind string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logging.GetLogger(logrus.Fields{"package": "v1"})
		c.Set("method", "single")

		// Create a new msg
		raw, err := c.GetRawData()
		if err != nil {
			// Log and return 200
			action := "read-body"
			c.Set("action", action)
			c.Set("error", err.Error())
			log.WithFields(logrus.Fields{"action": action}).Error(err.Error())
			c.AbortWithStatusJSON(http.StatusOK, returnJSON)
			return
		}

		// Unmarshal data into a Message
		msg, err := v1types.NewMessage(kind, raw)
		if err != nil {
			// Log and return 200
			action := "unmarshal"
			c.Set("action", action)
			c.Set("error", err.Error())
			log.WithFields(logrus.Fields{"action": action}).Error(err.Error())
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
		h.ingestionHandler.Write(msg)

		// Set additional metric data
		c.Set("event_count", 1)

		// Always return 200
		c.Header("Connection", "keep-alive")
		c.AbortWithStatusJSON(http.StatusOK, returnJSON)
	}
}
