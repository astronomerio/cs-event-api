package cmd

import (
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion"
	"github.com/spf13/cobra"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/api"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/config"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/logging"
	"github.com/sirupsen/logrus"
)

func buildAndStart() {
	apiServer := api.NewServer()
	appConfig := config.Get()
	appConfig.Print()

	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "cmd", "function": "main"})


	apiServerConfig := &api.ServerConfig{
		APIPort:          appConfig.APIPort,
		AdminPort:        appConfig.AdminPort,
		IngestionHandler: ingestion.NewHandler(appConfig.IngestionHandler),

		GracefulShutdownDelay: appConfig.GracefulShutdownDelay,
	}

	apiServer.WithConfig(apiServerConfig)
	apiServer.WithDefaultRoutes()

	apiServer.WithRequestID()

	if appConfig.HealthCheckEnabled {
		apiServer.WithHealthCheck()
	}

	if appConfig.PrometheusEnabled {
		apiServer.WithPrometheusMonitoring()
	}

	if appConfig.PProfEnabled {
		apiServer.WithPProf()
	}

	logger.Info("starting api server")
	apiServer.Run()
}

var RootCmd = &cobra.Command{
	Use: "clickstream-api",
	Run: func(cmd *cobra.Command, args []string) {
		buildAndStart()
	},
}
