package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

// Configuration is a stuct to hold event-api configs
type Configuration struct {
	APIPort               string
	AdminPort             string
	APIInterface          string
	AdminInterface        string
	GracefulShutdownDelay int
	MessageWriter         string
	StreamName            string
	KafkaTopic            string
	KafkaBrokers          []string
	PrometheusEnabled     bool
	HealthCheckEnabled    bool
	PProfEnabled          bool
	DebugMode             bool
	LogFormat             string
}

// AppConfig is a global instance of Configuration
var AppConfig Configuration

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	AppConfig = Configuration{}

	setDefaults()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed reading config file: %s\n", err)
	}

	viper.SetEnvPrefix("EA")
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

func setDefaults() {
	viper.SetDefault("PProfEnabled", false)
	viper.SetDefault("PrometheusEnabled", true)
	viper.SetDefault("HealthCheckEnabled", true)
	viper.SetDefault("GracefulShutdownDelay", 10)
	viper.SetDefault("APIPort", "8080")
	viper.SetDefault("AdminPort", "8081")
	viper.SetDefault("APIInterface", "0.0.0.0")
	viper.SetDefault("AdminInterface", "0.0.0.0")
	viper.SetDefault("DebugMode", false)
	viper.SetDefault("LogFormat", "json")
	viper.SetDefault("MessageWriter", "")
}

// Get returns the config
func Get() *Configuration {
	return &AppConfig
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
