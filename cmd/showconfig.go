package cmd

import (
	"fmt"

	"github.com/LandonTClipp/makemocks/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// NewShowConfigCmd returns the cobra command for generate
func NewShowConfigCmd(v *viper.Viper) *cobra.Command {
	return &cobra.Command{
		Use:   "showconfig",
		Short: "Print the final marshalled configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, _, err := GetGoodies(v)
			if err != nil {
				internal.StackAndFail(err)
			}
			out, err := yaml.Marshal(config)
			if err != nil {
				internal.StackAndFail(err)
			}
			fmt.Printf(string(out))
			return nil
		},
	}
}
