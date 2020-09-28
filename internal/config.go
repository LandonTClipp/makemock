package internal

import "github.com/spf13/viper"

// Version is the version of makemock
var Version = "0.0.0-dev"

// Config contains the configuration for the makemock application
type Config struct {
	// Config is the config file used, if any
	Config string `mapstructure:"config"`
	// DisableColor disables coloring of log messages
	DisableColor bool `mapstructure:"disable-color"`
	// LogLevel is the level of the logger
	LogLevel string `mapstructure:"log-level"`
}

// GetConfigDefault returns a Config struct with the application
// default values applied.
func GetConfigDefault() *Config {
	return &Config{
		Config:       "",
		DisableColor: false,
		LogLevel:     "info",
	}
}

// GetConfigFromViper unmarshals the viper object into a config.
// Any key in the Config struct that is not provided by viper
// will be initialized to a default value.
func GetConfigFromViper(v *viper.Viper) (*Config, error) {
	c := GetConfigDefault()
	return c, viper.UnmarshalExact(c)
}
