package logging

import (
	"github.com/astronomerio/clickstream-ingestion-api/pkg/config"
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
	}
}

func GetLogger() *logrus.Logger {
	return log
}
