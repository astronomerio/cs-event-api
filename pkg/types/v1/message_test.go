package v1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessage_String(t *testing.T) {
	staticTime := time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)
	m := Message{
		AppID:             "APP_ID",
		Timestamp:         GenericTime{staticTime},
		OriginalTimestamp: GenericTime{staticTime},
		ReceivedAt:        staticTime,
		SentAt:            GenericTime{staticTime},
	}

	expectedString := "{\"appId\":\"APP_ID\",\"timestamp\":\"2017-01-01T00:00:00Z\",\"originalTimestamp\":\"2017-01-01T00:00:00Z\",\"sentAt\":\"2017-01-01T00:00:00Z\",\"receivedAt\":\"2017-01-01T00:00:00Z\"}"

	assert.Equal(t, expectedString, m.String())
}
