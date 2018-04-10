package cmd

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/astronomerio/event-api/api"
	"github.com/astronomerio/event-api/api/v1"
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
	config.AppConfig.Print()

	// Create a waitgroup to ensure a clean shutdown.
	var wg sync.WaitGroup

	// Listen for system signals to shutdown and close our shutdown channel
	shutdownChan := make(chan struct{})
	go func() {
		sc := make(chan os.Signal)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)
		<-sc
		log.Info("Initiating shutdown sequence")
		close(shutdownChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Create a server config
		apiServerConfig := &api.ServerConfig{
			APIPort:               config.AppConfig.APIPort,
			AdminPort:             config.AppConfig.AdminPort,
			GracefulShutdownDelay: config.AppConfig.GracefulShutdownDelay,
		}

		// Create the producer
		producer := ingestion.NewMessageWriter(config.AppConfig.MessageWriter)
		defer producer.Close()

		// Create main server object
		apiServer := api.NewServer(apiServerConfig).
			WithRequestID().
			WithRouteHandler(v1.NewRouteHandler(producer))
		defer apiServer.Close()

		if config.AppConfig.PrometheusEnabled {
			apiServer.WithPrometheus()
		}

		if config.AppConfig.PProfEnabled {
			apiServer.WithPProf()
		}

		apiServer.Serve(shutdownChan)
	}()

	wg.Wait()
	log.Info("Finished")
}
