package v1

import (
	"encoding/json"
)

type Batch struct {
	SentAt   GenericTime `json:"sentAt,omitempty"`
	Messages []*Message  `json:"batch,omitempty"`
}

func (b *Batch) String() string {
	s, err := json.Marshal(b)
	if err != nil {
		return "<nil>"
	}
	return string(s)
}
