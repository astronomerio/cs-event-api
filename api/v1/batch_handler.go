package v1

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/arizz96/event-api/logging"
	v1types "github.com/arizz96/event-api/types/v1"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *RouteHandler) batchHandler(c *gin.Context) {
	log := logging.GetLogger(logrus.Fields{"package": "v1"})
	c.Set("method", "batch")

	// Read the raw bytes from the request
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

	// If gzipped, unzip and reset raw to unzipped data
	if c.GetHeader("Content-Encoding") == "gzip" {
		raw, err = unzip(raw)
		if err != nil {
			// Log and return 200
			action := "gzip-inflate"
			c.Set("action", action)
			c.Set("error", err.Error())
			log.WithFields(logrus.Fields{"action": action}).Error(err.Error())
			c.AbortWithStatusJSON(http.StatusOK, returnJSON)
			return
		}
	}

	// Create a batch object
	var batch v1types.Batch

	// Unmarshal data into a Batch
	err = json.Unmarshal(raw, &batch)
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

	// Loop over this batches messages
	for _, msg := range batch.Messages {

		// Apply batch level SentAt to each msg
		msg.WithSentAt(batch.SentAt)

		// Apply ReceivedAt date
		msg.WithReceivedAt(time.Now().UTC())

		// Skew timestamp to account for bad client clocks
		msg.SkewTimestamp()

		// Merge batch level context to msg context
		if batch.Context != nil && msg.GetContext() == nil {
			err := msg.MergeContext(batch.Context)
			if err != nil {
				action := "merge-context"
				c.Set("action", action)
				c.Set("error", err.Error())
				log.WithFields(logrus.Fields{"action": action}).Error(err.Error())
			}
		}

		// Merge batch level integrations to msg integrations
		if batch.Integrations != nil && msg.GetIntegrations() == nil {
			err = msg.MergeIntegrations(batch.Integrations)
			if err != nil {
				action := "merge-integrations"
				c.Set("action", action)
				c.Set("error", err.Error())
				log.WithFields(logrus.Fields{"action": action}).Error(err.Error())
			}
		}

		// Apply metadata from context
		msg.WithRequestMetadata(metadata)

		// Pass the msg along to the adapter
		h.ingestionHandler.Write(msg)
	}

	// Set additional metric data
	c.Set("event_count", len(batch.Messages))

	// Always return 200
	c.AbortWithStatusJSON(http.StatusOK, returnJSON)
}

// Unzip will unzip a gzipped payload and return the raw data
func unzip(b []byte) (data []byte, err error) {
	gzData, err := gzip.NewReader(bytes.NewBuffer(b))
	defer gzData.Close()

	if err != nil {
		return nil, err
	}

	d, err := ioutil.ReadAll(gzData)
	if err != nil {
		return nil, err
	}

	return d, nil
}
