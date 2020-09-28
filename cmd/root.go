package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/LandonTClipp/makemocks/internal"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "makemocks",
	Short:   "Generate mock objects using testify",
	Version: internal.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	pflags := rootCmd.PersistentFlags()
	pflags.StringVar(&cfgFile, "config", "", "config file (default is ./.makemock.yaml)")
	pflags.Bool("disable-color", false, "Disable coloring of log output")
	pflags.StringP("log-level", "l", "info", "Log level. Choose from: debug, info, warn, error, fatal.")
	viper.BindPFlags(pflags)

	v := viper.GetViper()
	rootCmd.AddCommand(NewGenerateCmd(v), NewShowConfigCmd(v))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".makemock" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName(".makemocks")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("MAKEMOCK")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// GetGoodies returns the application goodies (namely, configuration and context).
// Any return arguments besides error itself are undefined (i.e. possibly nill)
// on errors.
func GetGoodies(v *viper.Viper) (*internal.Config, context.Context, error) {
	config, err := internal.GetConfigFromViper(v)
	if err != nil {
		return config, nil, err
	}

	log, err := internal.GetNewLogger(config.DisableColor, config.LogLevel)
	if err != nil {
		return config, nil, err
	}
	ctx := log.WithContext(context.Background())
	return config, ctx, nil
}
