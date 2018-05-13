package types

import (
	"github.com/sirupsen/logrus"
	"github.com/arizz96/event-api/logging"
)

type ConvertibleBoolean bool

func (bit *ConvertibleBoolean) UnmarshalJSON(data []byte) error {
	log := logging.GetLogger(logrus.Fields{"package": "types"})

  asString := string(data)
	if asString == "1" || asString == "true" {
		*bit = true
	} else if asString == "0" || asString == "false" {
		*bit = false
	} else {
		log.WithFields(logrus.Fields{"action": "Error during ConvertibleBoolean unmarshal"})
	}

	return nil
}
