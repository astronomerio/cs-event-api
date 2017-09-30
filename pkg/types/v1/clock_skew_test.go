package v1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ApplyClockSkew(t *testing.T) {
	tests := []struct {
		name   string
		input  *Message
		output *Message
	}{
		{
			"doesnt apply if Timestamp isnt present",
			&Message{},
			&Message{},
		},
		{
			"doesnt apply if SentAt isnt present",
			&Message{},
			&Message{},
		},
		{
			"doesnt apply if ReceivedAt isnt present",
			&Message{},
			&Message{},
		},
		{
			// should add 1.3s to the timestamp
			"correctly applies skew by correct amount",
			&Message{
				ReceivedAt: time.Date(2017, 1, 1, 0, 0, 2, 300*1e6, time.UTC),
				SentAt:     GenericTime{time.Date(2017, 1, 1, 0, 0, 1, 0, time.UTC)},
				Timestamp:  GenericTime{time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)},
				// OriginalTimestamp will get added by the function
			},
			&Message{
				Timestamp:         GenericTime{time.Date(2017, 1, 1, 0, 0, 1, 300*1e6, time.UTC)},
				OriginalTimestamp: GenericTime{time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origTimestamp := tt.input.Timestamp
			tt.input.SkewTimestamp()
			assert.Equal(t, tt.input.Timestamp, tt.output.Timestamp)
			assert.Equal(t, origTimestamp, tt.output.OriginalTimestamp)
		})
	}
}
