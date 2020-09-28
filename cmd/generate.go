package cmd

import (
	"context"
	"strings"

	"github.com/LandonTClipp/makemocks/internal"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/tools/go/packages"
)

// NewGenerateCmd returns the cobra command for generate
func NewGenerateCmd(v *viper.Viper) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate mocks",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, ctx, err := GetGoodies(v)
			if err != nil {
				internal.StackAndFail(err)
			}
			generate, err := GetGenerateFromConfig(&GenerateConfig{
				Packages: config.Packages,
			})
			if err != nil {
				internal.StackAndFail(err)
			}
			if err := generate.Run(ctx); err != nil {
				internal.StackAndFail(err)
			}
			return nil
		},
	}
}

type GenerateConfig struct {
	Packages map[string]internal.Package
}

// Generate is the application object for the generate command
type Generate struct {
	Config *GenerateConfig
}

// GetGenerateFromConfig returns a Generate object using the provided configuration object
func GetGenerateFromConfig(c *GenerateConfig) (*Generate, error) {
	return &Generate{
		Config: c,
	}, nil
}

// Run runs the command
func (g *Generate) Run(ctx context.Context) error {
	log := zerolog.Ctx(ctx).With().Str(internal.LogKeyCommand, "generate").Logger()
	ctx = log.WithContext(ctx)

	log.Debug().Msgf("%+v", g.Config)
	packageNames := []string{}
	for name, val := range g.Config.Packages {
		log.Debug().Msgf(name)
		log.Debug().Msgf(val.Test1)
		log.Debug().Msgf(val.Test2)
		if name == internal.PackageDefault {
			continue
		}
		packageNames = append(packageNames, name)
	}
	log.Info().Msgf("loading packages: %s", strings.Join(packageNames, ", "))
	foundPackages, err := packages.Load(nil, packageNames...)
	if err != nil {
		return errors.Wrapf(err, "failed to load packages")
	}
	for _, p := range foundPackages {
		log.Info().Msgf(p.PkgPath)
	}
	return nil
}
