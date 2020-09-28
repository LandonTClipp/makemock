package cmd

import (
	"context"

	"github.com/LandonTClipp/makemock/internal"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate mocks",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, ctx, err := GetGoodies(viper.GetViper())
		if err != nil {
			internal.StackAndFail(err)
		}
		generate, err := GetGenerateFromConfig(config)
		if err != nil {
			internal.StackAndFail(err)
		}
		if err := generate.Run(ctx); err != nil {
			internal.StackAndFail(err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

// Generate is the application object for the generate command
type Generate struct {
	Config *internal.Config
}

// GetGenerateFromConfig returns a Generate object using the provided configuration object
func GetGenerateFromConfig(c *internal.Config) (*Generate, error) {
	return &Generate{
		Config: c,
	}, nil
}

// Run runs the command
func (g *Generate) Run(ctx context.Context) error {
	log := zerolog.Ctx(ctx)
	log.Debug().Msgf("hello")
	log.Info().Msgf("hello")
	log.Warn().Msgf("Hello")
	log.Error().Msgf("Hello")
	log.Fatal().Msgf("Hello")
	return nil
}
