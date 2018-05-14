package v1

import (
	"encoding/base64"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequestMetadata is data pulled or generated from outside the message body
type RequestMetadata struct {
	IP       net.IP
	WriteKey string
}

// NewRequestMetadata generates a RequestMetadata for a request
func NewRequestMetadata(c *gin.Context) (md RequestMetadata) {
	// Assign IP address
	md.IP = net.ParseIP(c.ClientIP())

	// Grab writeKey from auth header
	authHeader := strings.TrimLeft(c.GetHeader("Authorization"), "Basic ")
	if authHeader != "" {
		bs, err := base64.StdEncoding.DecodeString(authHeader)
		if err != nil {
			return
		}
		md.WriteKey = strings.TrimRight(string(bs), ":")
	}
	return
}
