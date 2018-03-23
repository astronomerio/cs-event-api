package v1

import "encoding/json"

// NewMessage returns a new message of the given type
func NewMessage(kind string, raw []byte) (msg Message, err error) {
	switch kind {
	case "alias":
		msg = new(Alias)
	case "group":
		msg = new(Group)
	case "identify":
		msg = new(Identify)
	case "page":
		msg = new(Page)
	case "screen":
		msg = new(Screen)
	case "track":
		msg = new(Track)
	}
	err = json.Unmarshal(raw, &msg)
	return
}
