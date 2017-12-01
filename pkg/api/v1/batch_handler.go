package v1

import (
	"encoding/json"
	"log"
	"net/http"

	v1types "github.com/astronomerio/cs-event-api/pkg/types/v1"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *RouteHandler) batchHandler(c *gin.Context) {
	logger := h.logger.WithFields(logrus.Fields{"package": "v1", "function": "batchHandler"})
	c.Set("profile", true)
	c.Set("type", "batch")
	c.Set("action", "batch")

	rd, err := c.GetRawData()

	if err != nil {
		logger.Error(logrus.Fields{"stage": "1", "error": err.Error()})
		c.Set("stage", "1")
		log.Println(err.Error())
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
