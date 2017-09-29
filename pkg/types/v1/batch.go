package v1

import (
	"encoding/json"
	"time"
)

type Batch struct {
	SentAt   time.Time `json:"sentAt,omitempty"`
	Messages []Message `json:"batch,omitempty"`
}

func (b *Batch) String() string {
	s, err := json.Marshal(b)
	if err != nil {
		return "<nil>"
	}
	return string(s)
}
