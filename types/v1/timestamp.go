package v1

import (
	"bytes"
	"time"

	"github.com/araddon/dateparse"
)

// Format is ISO-8601 format
const Format = "2006-01-02T15:04:05"
const jsonFormat = `"` + Format + `"`

// Timestamp represents an ISO-8601 date string
type Timestamp struct {
	time.Time
}

// MarshalJSON outputs an ISO-8601 date string
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t.Time).Format(jsonFormat)), nil
}

// UnmarshalJSON unmarshals an ISO-8601 date string
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	ts, err := dateparse.ParseAny(string(bytes.Trim(data, "\"")))
	if err != nil {
		return err
	}
	t.Time = ts
	return nil
}

func (t Timestamp) String() string {
	return t.Time.String()
}
