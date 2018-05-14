package v1

import "encoding/json"

// Track represents a track call.
type Track struct {
	Event
	EventName  string     `json:"event"`
	Properties Properties `json:"properties,omitempty"`
}

// Validate validates the message
func (msg Track) Validate() error {
	if len(msg.EventName) == 0 {
		return FieldError{
			Type:  "Track",
			Name:  "Event",
			Value: msg.EventName,
		}
	}

	if len(msg.UserID) == 0 && len(msg.AnonymousID) == 0 {
		return FieldError{
			Type:  "Track",
			Name:  "UserID",
			Value: msg.UserID,
		}
	}

	return nil
}

func (msg Track) String() string {
	bs, err := json.Marshal(msg)
	if err != nil {
		return ""
	}
	return string(bs)
}
