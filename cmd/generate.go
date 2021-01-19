package cmd

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
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

type parserEntry struct {
	fileName string
	pkg      *packages.Package
	syntax   *ast.File
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
	foundPackages, err := packages.Load(&packages.Config{
		Mode: packages.NeedFiles | packages.NeedSyntax,
	}, packageNames...)
	if err != nil {
		return errors.Wrapf(err, "failed to load packages")
	}

	parsers := []*ast.File{}

	for _, p := range foundPackages {
		if len(p.Errors) > 0 {
			return p.Errors[0]
		}
		log.Info().Msgf(p.PkgPath)
		for _, f := range p.GoFiles {
			log.Info().Msgf(f)
			fset := token.NewFileSet()
			fileParser, err := parser.ParseFile(fset, f, nil, parser.ParseComments)
			if err != nil {
				return errors.Wrapf(err, "failed to parser file")
			}
			parsers = append(parsers, fileParser)
		}
	}

	for _, fileParser := range parsers {
		fileLog := log.With().Str("file", fileParser.Name.Name).Logger()
		fileLog.Info().Msgf("printing declared functions")

		for _, f := range fileParser.Decls {
			fn, ok := f.(*ast.FuncDecl)
			if !ok {
				continue
			}
			fileLog.Info().Msgf(fn.Name.Name)
		}

		fileLog.Info().Msgf("printing declared interfaces")

		for _, node := range parsers {
			ast.Inspect(node, func(n ast.Node) bool {
				switch t := n.(type) {
				case *ast.TypeSpec:
					switch t.Type.(type) {
					case *ast.InterfaceType:
						log.Info().Msgf(t.Name.Name)
					}
				}
				return true
			})

		}
	}
	return nil
}
