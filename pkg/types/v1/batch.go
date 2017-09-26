package v1

import "time"

type Batch struct {
	SentAt   time.Time `json:"sentAt,omitempty"`
	Messages []Message `json:"batch,omitempty"`
}
