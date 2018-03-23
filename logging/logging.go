package logging

import (
	"github.com/astronomerio/event-api/config"
	"github.com/sirupsen/logrus"
)

// Singleton logger for application
var log *logrus.Logger

// Configure logger on startup
func init() {
	c := config.Get()
	log = logrus.New()

	if c.LogFormat == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	if c.DebugMode {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
}

// GetLogger returns the singleton logger
func GetLogger() *logrus.Logger {
	return log
}
