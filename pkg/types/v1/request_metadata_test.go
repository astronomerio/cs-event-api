package v1

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRequestMetadata(t *testing.T) {
	ip := "8.8.8.8"
	appID := "APP_ID"
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()

	t.Run("will return the correct IP and not decode auth header because its not present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(rec)

		c.Request, _ = http.NewRequest("GET", "/v1/t", nil)
		c.Request.RemoteAddr = ip + ":80"

		md1 := GetRequestMetadata(c)
		assert.Equal(t, ip, md1.IP)
		assert.Equal(t, "", md1.AppID)
	})

	t.Run("will decode the authorization header", func(t *testing.T) {
		c, _ := gin.CreateTestContext(rec)

		c.Request, _ = http.NewRequest("GET", "/v1/t", nil)

		c.Request.RemoteAddr = ip + ":80"
		encodedID := base64.StdEncoding.EncodeToString([]byte(appID))
		c.Request.Header.Set("authorization", encodedID)

		md1 := GetRequestMetadata(c)
		assert.Equal(t, ip, md1.IP)
		assert.Equal(t, appID, md1.AppID)
	})
}
