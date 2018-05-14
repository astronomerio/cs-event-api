package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

// AppConfig is the gloabl application configuration
var AppConfig = &Configuration{}

// Configuration is a stuct to hold event-api configs
type Configuration struct {
	DebugMode             bool   `mapstructure:"DEBUG_MODE"`
	LogFormat             string `mapstructure:"LOG_FORMAT"`
	APIPort               string `mapstructure:"API_PORT"`
	AdminPort             string `mapstructure:"ADMIN_PORT"`
	APIInterface          string `mapstructure:"API_INTERFACE"`
	AdminInterface        string `mapstructure:"ADMIN_INTERFACE"`
	GracefulShutdownDelay int    `mapstructure:"GRACEFUL_SHUTDOWN_DELAY"`
	MessageWriter         string `mapstructure:"MESSAGE_WRITER"`
	KafkaTopic            string `mapstructure:"KAFKA_TOPIC"`
	KafkaBrokers          string `mapstructure:"KAFKA_BROKERS"`
	PrometheusEnabled     bool   `mapstructure:"PROMETHEUS_ENABLED"`
	HealthCheckEnabled    bool   `mapstructure:"HEALTHCHECK_ENABLED"`
	PProfEnabled          bool   `mapstructure:"PPROF_ENABLED"`
	FlushTimeout          int    `mapstructure:"FLUSH_TIMEOUT"`
	QueueBufferingDelayMs int    `mapstructure:"QUEUE_BUFFERING_DELAY_MS"`
}

func init() {
	appViper := viper.New()
	appViper.SetEnvPrefix("EA")
	appViper.AutomaticEnv()

	appViper.SetDefault("DEBUG_MODE", false)
	appViper.SetDefault("LOG_FORMAT", "json")
	appViper.SetDefault("API_PORT", "8080")
	appViper.SetDefault("ADMIN_PORT", "8081")
	appViper.SetDefault("API_INTERFACE", "0.0.0.0")
	appViper.SetDefault("ADMIN_INTERFACE", "0.0.0.0")
	appViper.SetDefault("PROMETHEUS_ENABLED", true)
	appViper.SetDefault("HEALTHCHECK_ENABLED", true)
	appViper.SetDefault("GRACEFUL_SHUTDOWN_DELAY", 10)
	appViper.SetDefault("PPROF_ENABLED", false)
	appViper.SetDefault("MESSAGE_WRITER", "")
	appViper.SetDefault("KAFKA_BROKERS", "")
	appViper.SetDefault("KAFKA_TOPIC", "")
	appViper.SetDefault("FLUSH_TIMEOUT", 10000)
	appViper.SetDefault("QUEUE_BUFFERING_DELAY_MS", 5000)

	if err := appViper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}

// Print prints the configuration to stdout
func (c *Configuration) Print() {
	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	fmt.Println("=============== Configuration ===============")
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fmt.Printf("%s %s = %v\n", t.Field(i).Name, f.Type(), f.Interface())
	}
	fmt.Println("=============================================")
}
