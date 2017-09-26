package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	APIPort               string
	AdminPort             string
	GracefulShutdownDelay int

	IngestionHandler string

	StreamName string

	KafkaTopic   string
	KafkaBrokers []string

	PrometheusEnabled  bool
	HealthCheckEnabled bool
	PProfEnabled       bool
}

var AppConfig Configuration

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	setDefaults()

	AppConfig = Configuration{}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

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
