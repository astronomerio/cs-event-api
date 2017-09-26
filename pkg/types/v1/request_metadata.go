package v1

import (
	"encoding/base64"

	"github.com/gin-gonic/gin"
)

type RequestMetadata struct {
	IP    string
	AppID string
}

func GetRequestMetadata(c *gin.Context) (md RequestMetadata) {
	md.IP = c.ClientIP()

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
