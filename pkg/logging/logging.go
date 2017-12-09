package logging

import (
	"github.com/astronomerio/event-api/pkg/config"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	c := config.Get()
	if c.LogFormat == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	log = logrus.New()
	if c.DebugMode {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
}

func GetLogger() *logrus.Logger {
	return log
}
