package v1

import "encoding/json"

// Group represents a group call.
type Group struct {
	Event
	GroupID string `json:"groupId"`
}

func (msg Group) validate() error {
	if len(msg.GroupID) == 0 {
		return FieldError{
			Type:  "Group",
			Name:  "GroupID",
			Value: msg.GroupID,
		}
	}

	if len(msg.UserID) == 0 && len(msg.AnonymousID) == 0 {
		return FieldError{
			Type:  "Group",
			Name:  "UserID",
			Value: msg.UserID,
		}
	}

	return nil
}

func (msg Group) String() string {
	bs, err := json.Marshal(msg)
	if err != nil {
		return ""
	}
	return string(bs)
}
