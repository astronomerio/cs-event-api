package v1

import (
	"log"
	"net/http"
	"fmt"

	v1types "github.com/astronomerio/clickstream-ingestion-api/pkg/types/v1"
	"github.com/gin-gonic/gin"
)

func (h *RouteHandler) importHandler(c *gin.Context) {
	c.Set("profile", true)
	c.Set("type", "import")
	c.Set("action", "import")

	rd, err := c.GetRawData()

	if err != nil {
		c.Set("error", err.Error())
		c.Set("stage", "1")
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusOK, returnJSON)
		return
	}

	batch, err := gzipToBatch(rd)
	if err != nil {
		c.Set("error", err.Error())
		c.Set("stage", "2")
		c.AbortWithStatusJSON(http.StatusOK, returnJSON)
		return
	}
	fmt.Println("num messages:", len(batch.Messages))

	md := v1types.GetRequestMetadata(c)
	for _, m := range batch.Messages {
		m.SentAt = batch.SentAt
		m.ApplyMetadata(md)
		m.SkewTimestamp()
		h.ingestionHandler.ProcessMessage(m.String(), m.PartitionKey())
	}

	c.AbortWithStatusJSON(http.StatusOK, returnJSON)
}
