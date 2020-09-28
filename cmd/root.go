package cmd

import (
	"context"
	"fmt"

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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
}

var v *viper.Viper

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(
		NewGenerateCmd(v),
		NewShowConfigCmd(v),
	)
	if err := rootCmd.Execute(); err != nil {
		internal.StackAndFail(err)
	}
}

func init() {
	v = viper.NewWithOptions(viper.KeyDelimiter("rf66bdg4ot554a00lil2k8o2ynjhi"))

	cobra.OnInitialize(initConfig)

	pflags := rootCmd.PersistentFlags()
	pflags.StringVar(&cfgFile, "config", "", "config file (default is ./.makemocks.yaml)")
	pflags.Bool("disable-color", false, "Disable coloring of log output")
	pflags.StringP("log-level", "l", "info", "Log level. Choose from: debug, info, warn, error, fatal.")
	v.BindPFlags(pflags)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".makemocks" (without extension).
		v.AddConfigPath(".")
		v.SetConfigName(".makemocks")
	}

	v.AutomaticEnv() // read in environment variables that match
	v.SetEnvPrefix("MAKEMOCKS")

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", v.ConfigFileUsed())
	}
}

// GetGoodies returns the application goodies (namely, configuration and context).
// Any return arguments besides error itself are undefined (i.e. possibly nil)
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
