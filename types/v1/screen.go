package v1

import "encoding/json"

// Screen represents a screen call.
type Screen struct {
	Event
	Name       string     `json:"name,omitempty"`
	Properties Properties `json:"properties,omitempty"`
}

func (msg Screen) validate() error {
	if len(msg.UserID) == 0 && len(msg.AnonymousID) == 0 {
		return FieldError{
			Type:  "Screen",
			Name:  "UserID",
			Value: msg.UserID,
		}
	}

	return nil
}

func (msg Screen) String() string {
	bs, err := json.Marshal(msg)
	if err != nil {
		return ""
	}
	return string(bs)
}
