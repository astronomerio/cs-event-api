package v1

import (
	"bytes"
	"fmt"
	"time"

	"github.com/araddon/dateparse"
)

type GenericTime struct {
	time.Time
}

func (gt *GenericTime) UnmarshalJSON(d []byte) error {
	t, err := dateparse.ParseAny(string(bytes.Trim(d, "\"")))
	if err != nil {
		return fmt.Errorf("GenericTime.UnmarshalJSON: %s", err.Error())
	}
	gt.Time = t
	return nil
}
