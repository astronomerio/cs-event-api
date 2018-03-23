package v1

import "encoding/json"

// Alias represents an alias call
type Alias struct {
	Event
	PreviousID string `json:"previousId"`
}

func (msg Alias) validate() error {
	if len(msg.UserID) == 0 {
		return FieldError{
			Type:  "Alias",
			Name:  "UserID",
			Value: msg.UserID,
		}
	}

	if len(msg.PreviousID) == 0 {
		return FieldError{
			Type:  "Alias",
			Name:  "PreviousID",
			Value: msg.PreviousID,
		}
	}

	return nil
}

func (msg Alias) String() string {
	bs, err := json.Marshal(msg)
	if err != nil {
		return ""
	}
	return string(bs)
}
