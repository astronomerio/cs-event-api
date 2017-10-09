package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	APIPort               string
	AdminPort             string
	APIInterface          string
	AdminInterface        string
	GracefulShutdownDelay int

	IngestionHandler string

	StreamName string

	KafkaTopic   string
	KafkaBrokers []string

	PrometheusEnabled  bool
	HealthCheckEnabled bool
	PProfEnabled       bool

	DebugMode bool

	LogFormat string
}

var AppConfig Configuration

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if sandboxPath, ok := os.LookupEnv("MESOS_SANDBOX"); ok {
		viper.AddConfigPath(sandboxPath)
	}

	setDefaults()

	AppConfig = Configuration{}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.SetEnvPrefix("ingestion")
	viper.AutomaticEnv()

	viper.Set("AdminPort", viper.GetString("admin_port"))
	viper.Set("APIPort", viper.GetString("api_port"))
	viper.Set("PrometheusEnabled", viper.GetBool("enable_prometheus"))
	viper.Set("DebugMode", viper.GetBool("debug_mode"))
	viper.Set("AdminInterface", viper.GetString("admin_interface"))
	viper.Set("HealthCheckEnabled", viper.GetBool("enable_health_check"))
	viper.Set("GracefulShutdownDelay", viper.GetInt("graceful_shutdown"))
	viper.Set("APIInterface", viper.GetString("api_interface"))
	viper.Set("PProfEnabled", viper.GetBool("enable_pprof"))


	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

func setDefaults() {
	viper.SetDefault("PProfEnabled", false)
	viper.SetDefault("PrometheusEnabled", true)
	viper.SetDefault("HealthCheckEnabled", true)
	viper.SetDefault("GracefulShutdownDelay", 30)
	viper.SetDefault("APIPort", "8080")
	viper.SetDefault("AdminPort", "8081")
	viper.SetDefault("APIInterface", "0.0.0.0")
	viper.SetDefault("AdminInterface", "0.0.0.0")
	viper.SetDefault("DebugMode", false)
	viper.SetDefault("LogFormat", "json")
}

// Get returns the config
func Get() *Configuration {
	return &AppConfig
}

// Print prints the configuration to stdout
func (c *Configuration) Print() {
	fmt.Println("=================")
	fmt.Printf("IngestionHandler: %s\n", c.IngestionHandler)
	fmt.Printf("StreamName: %s\n", c.StreamName)
	fmt.Printf("KafkaTopic: %s\n", c.KafkaTopic)
	fmt.Printf("KafkaBrokers: %s\n", c.KafkaBrokers)
	fmt.Printf("PrometheusEnabled: %t\n", c.PrometheusEnabled)
	fmt.Printf("HealthCheckEnabled: %t\n", c.HealthCheckEnabled)
	fmt.Printf("PProfEnabled: %t\n", c.PProfEnabled)
	fmt.Printf("APIPort: %s\n", c.APIPort)
	fmt.Printf("AdminPort: %s\n", c.AdminPort)
	fmt.Printf("GracefulShutdownDelay: %d\n", c.GracefulShutdownDelay)
	fmt.Println("=================")
}
