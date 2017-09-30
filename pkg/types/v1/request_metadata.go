package v1

import (
	"encoding/base64"
	"time"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/util"
	"github.com/gin-gonic/gin"
)

type RequestMetadata struct {
	IP         string
	AppID      string
	ReceivedAt time.Time
}

func GetRequestMetadata(c *gin.Context) (md RequestMetadata) {
	md.IP = c.ClientIP()
	md.ReceivedAt = util.NowUTC()

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		bs, err := base64.StdEncoding.DecodeString(authHeader)
		if err != nil {
			// TODO: handle error
			return
		}
		md.AppID = string(bs)
	}
	return
}
