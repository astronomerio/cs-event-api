package v1

import "encoding/json"

// Identify represents an identify call.
type Identify struct {
	Event
	Traits Traits `json:"traits,omitempty"`
}

func (msg Identify) validate() error {
	if len(msg.UserID) == 0 && len(msg.AnonymousID) == 0 {
		return FieldError{
			Type:  "Identify",
			Name:  "UserID",
			Value: msg.UserID,
		}
	}

	return nil
}

func (msg Identify) String() string {
	bs, err := json.Marshal(msg)
	if err != nil {
		return ""
	}
	return string(bs)
}
