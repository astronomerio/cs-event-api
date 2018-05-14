package v1

import "encoding/json"

// Page represents a page call.
type Page struct {
	Event
	Name       string     `json:"name,omitempty"`
	Properties Properties `json:"properties,omitempty"`
}

// Validate validates the message
func (msg Page) Validate() error {
	if len(msg.UserID) == 0 && len(msg.AnonymousID) == 0 {
		return FieldError{
			Type:  "Page",
			Name:  "UserID",
			Value: msg.UserID,
		}
	}

	return nil
}

func (msg Page) String() string {
	bs, err := json.Marshal(msg)
	if err != nil {
		return ""
	}
	return string(bs)
}
