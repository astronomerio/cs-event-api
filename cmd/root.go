package cmd

import (
	"github.com/astronomerio/clickstream-ingestion-api/pkg/ingestion"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/logger"
	"github.com/spf13/cobra"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/api"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/config"
)

func buildAndStart() {
	apiserver := api.NewServer()
	appConfig := config.Get()
	appConfig.Print()

	apiserverConfig := &api.APIServerConfig{
		APIPort:          appConfig.APIPort,
		AdminPort:        appConfig.AdminPort,
		Logger:           logger.NewLogger("mock"),
		IngestionHandler: ingestion.NewHandler(appConfig.IngestionHandler),

		GracefulShutdownDelay: appConfig.GracefulShutdownDelay,
	}

	apiserver.WithConfig(apiserverConfig)
	apiserver.WithDefaultRoutes()

	apiserver.WithRequestID()

	if appConfig.HealthCheckEnabled {
		apiserver.WithHealthCheck()
	}

	if appConfig.PrometheusEnabled {
		apiserver.WithPrometheusMonitoring()
	}

	if appConfig.PProfEnabled {
		apiserver.WithPProf()
	}

	apiserver.Run()
}

var RootCmd = &cobra.Command{
	Use: "clickstream-api",
	Run: func(cmd *cobra.Command, args []string) {
		buildAndStart()
	},
}
