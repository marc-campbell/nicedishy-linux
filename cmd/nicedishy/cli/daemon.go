package cli

import (
	"github.com/marc-campbell/nicedishy/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DaemonCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Starts the NiceDishy agent",
		Long:  `starts the NiceDishy agent`,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()

			if v.GetString("log-level") == "debug" {
				logger.Info("setting log level to debug")
				logger.SetDebug()
			}

			return nil
		},
	}

	cmd.Flags().String("log-level", "info", "set the log level")

	return cmd
}
