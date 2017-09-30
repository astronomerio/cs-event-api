package v1

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	v1types "github.com/astronomerio/clickstream-ingestion-api/pkg/types/v1"
	"github.com/gin-gonic/gin"
)

func gzipToBatch(b []byte) (batch v1types.Batch, err error) {
	gzData, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return
	}
	defer gzData.Close()
	d, err := ioutil.ReadAll(gzData)
	if err != nil {
		return
	}
	err = json.Unmarshal(d, &batch)
	return
}

func (h *RouteHandler) batchHandler(c *gin.Context) {
	c.Set("type", "batch")
	c.Set("action", "batch")
	log.Println("batch")

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

	md := v1types.GetRequestMetadata(c)
	for _, m := range batch.Messages {
		m.SentAt = batch.SentAt
		m.ApplyMetadata(md)
		m.SkewTimestamp()
		h.ingestionHandler.ProcessMessage(m.String(), m.PartitionKey())
	}

	fmt.Println("num messages:", len(batch.Messages))

	c.AbortWithStatusJSON(http.StatusOK, returnJSON)
}
