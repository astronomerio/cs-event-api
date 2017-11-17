package v1

import (
	"encoding/json"
)

type Batch struct {
	SentAt       GenericTime            `json:"sentAt,omitempty"`
	Messages     []*Message             `json:"batch,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
	Context      map[string]interface{} `json:"traits,omitempty"`
}

func (b *Batch) String() string {
	s, err := json.Marshal(b)
	if err != nil {
		return "<nil>"
	}
	return string(s)
}
