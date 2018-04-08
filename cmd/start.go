package cmd

import (
	"github.com/astronomerio/event-api/api"
	"github.com/astronomerio/event-api/config"
	"github.com/astronomerio/event-api/ingestion"
	"github.com/astronomerio/event-api/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Event API",
	Long:  "Start the Event API",
	Run:   start,
}

func init() {
	RootCmd.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) {
	log := logging.GetLogger(logrus.Fields{"package": "cmd"})

	// Create main server object
	apiServer := api.NewServer()

	// Grab and print application config
	config.AppConfig.Print()

	// Create a server config
	apiServerConfig := &api.ServerConfig{
		APIPort:               config.AppConfig.APIPort,
		AdminPort:             config.AppConfig.AdminPort,
		MessageWriter:         ingestion.NewMessageWriter(config.AppConfig.MessageWriter),
		GracefulShutdownDelay: config.AppConfig.GracefulShutdownDelay,
	}

	// Set up our server options
	apiServer.
		WithConfig(apiServerConfig).
		WithDefaultRoutes().
		WithRequestID()

	if config.AppConfig.HealthCheckEnabled {
		apiServer.WithHealthCheck()
	}

	if config.AppConfig.PrometheusEnabled {
		apiServer.WithPrometheusMonitoring()
	}

	if config.AppConfig.PProfEnabled {
		apiServer.WithPProf()
	}

	log.Info("Starting API server")
	apiServer.Run()
}
