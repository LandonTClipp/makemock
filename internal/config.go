package internal

import (
	"github.com/chigopher/pathlib"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Version is the version of makemock
var Version = "0.0.0-dev"

// Config contains the global configuration for the makemock application.
// It is intended to be used to contain the marshalled application data.
// Parameters in here will need to be passed to the command-specific config,
// depending on what configuration the command needs.
type Config struct {
	// Config is the config file used, if any
	Config string `mapstructure:"config"`
	// DisableColor disables coloring of log messages
	DisableColor bool `mapstructure:"disable-color"`
	// LogLevel is the level of the logger
	LogLevel string             `mapstructure:"log-level"`
	Packages map[string]Package `mapstructure:"packages"`
}

// GetConfigDefault returns a Config struct with the application
// default values applied.
func GetConfigDefault() *Config {
	return &Config{
		Config:       "",
		DisableColor: false,
		LogLevel:     "info",
		Packages:     make(map[string]Package),
	}
}

// GetConfigFromViper unmarshals the viper object into a config.
// Any key in the Config struct that is not provided by viper
// will be initialized to a default value.
func GetConfigFromViper(v *viper.Viper) (*Config, error) {
	c := GetConfigDefault()
	if err := v.UnmarshalExact(c); err != nil {
		return c, errors.Wrapf(err, "failed to unmarshal")
	}
	// Because of an issue with viper (https://github.com/spf13/viper/issues/1014),
	// configuration keys become case insensitive. We need to hack around the
	// `packages` key in order for the configuration here to be case-sensitive.
	// So, unmarshal the raw yaml and replace c.Packages with the result of that.
	configPath := pathlib.NewPath(v.ConfigFileUsed())
	configBytes, err := configPath.ReadFile()
	if err != nil {
		return c, errors.Wrapf(err, "failed to read config file")
	}
	cRaw := &Config{}
	if err := yaml.Unmarshal(configBytes, cRaw); err != nil {
		return c, errors.Wrapf(err, "failed to unmarshal config yaml")
	}
	c.Packages = cRaw.Packages
	return c, nil
}

type Package struct {
	Test1 string `mapstructure:"test1"`
	Test2 string `mapstructure:"test2"`
}

const (
	PackageDefault = "DEFAULT"
)
