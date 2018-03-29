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
	log := logging.GetLogger().WithFields(logrus.Fields{"package": "cmd"})

	// Create main server object
	apiServer := api.NewServer()

	// Grab and print application config
	appConfig := config.Get()
	appConfig.Print()

	// Create a server config
	apiServerConfig := &api.ServerConfig{
		APIPort:               appConfig.APIPort,
		AdminPort:             appConfig.AdminPort,
		MessageWriter:         ingestion.NewMessageWriter(appConfig.MessageWriter),
		GracefulShutdownDelay: appConfig.GracefulShutdownDelay,
	}

	// Set up our server options
	apiServer.
		WithConfig(apiServerConfig).
		WithDefaultRoutes().
		WithRequestID()

	if appConfig.HealthCheckEnabled {
		apiServer.WithHealthCheck()
	}

	if appConfig.PrometheusEnabled {
		apiServer.WithPrometheusMonitoring()
	}

	if appConfig.PProfEnabled {
		apiServer.WithPProf()
	}

	log.Info("Starting API server")
	apiServer.Run()
}
