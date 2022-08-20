package cli

import (
	"fmt"

	"github.com/marc-campbell/nicedishy-linux/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the version and exits",
		Long:  `Prints the version and exits`,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("NiceDishy Version %s\n", version.Version())

			return nil
		},
	}

	return cmd
}
