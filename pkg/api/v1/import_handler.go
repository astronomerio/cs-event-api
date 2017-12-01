package v1

import (
	"net/http"

	"encoding/json"

	v1types "github.com/astronomerio/cs-event-api/pkg/types/v1"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *RouteHandler) importHandler(c *gin.Context) {
	logger := h.logger.WithFields(logrus.Fields{"package": "v1", "function": "importHandler"})

	c.Set("profile", true)
	c.Set("type", "import")
	c.Set("action", "import")

	rd, err := c.GetRawData()

	if err != nil {
		h.logger.Error(logrus.Fields{"stage": "2", "error": err.Error()})
		c.Set("stage", "1")
		c.AbortWithStatusJSON(http.StatusOK, returnJSON)
		return
	}

	var batch v1types.Batch

	if c.GetHeader("Content-Encoding") == "gzip" {
		batch, err = gzipToBatch(rd)
		if err != nil {
			logger.Error(logrus.Fields{"stage": "2", "action": "gzip-inflate", "error": err.Error()})
			c.Set("error", err.Error())
			c.Set("stage", "2")
			c.AbortWithStatusJSON(http.StatusOK, returnJSON)
			return
		}
	} else {
		err = json.Unmarshal(rd, &batch)
		if err != nil {
			logger.Error(logrus.Fields{"stage": "2", "action": "batch-unmarshal", "error": err.Error()})
			c.Set("stage", "2")
			c.AbortWithStatusJSON(http.StatusOK, returnJSON)
			return
		}
	}

	md := v1types.GetRequestMetadata(c)
	for _, m := range batch.Messages {
		m.SentAt = batch.SentAt
		m.ApplyMetadata(md)
		m.SkewTimestamp()

		{
			err := mergeFields(&m.Context, batch.Context)
			if err != nil {
				logger.Error(logrus.Fields{"appID": m.AppID, "action": "merge-integrations", "error": err.Error()})
			}
		}
		{
			err := mergeFields(&m.Integrations, batch.Integrations)
			if err != nil {
				logger.Error(logrus.Fields{"appID": m.AppID, "action": "merge-integrations", "error": err.Error()})
			}
		}

		h.ingestionHandler.ProcessMessage(m.String(), m.PartitionKey())
	}

	c.AbortWithStatusJSON(http.StatusOK, returnJSON)
}
