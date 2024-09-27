package cmd

import (
	"github.com/mheers/vcluster-operator/helpers"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create [name]",
		Short: "creates a vcluster",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			client, err := initClient(true)
			if err != nil {
				return err
			}

			if len(args) == 0 {
				return cmd.Help()
			}
			name := args[0]
			cluster, err := client.Create(name)
			if err != nil {
				return err
			}

			return renderCluster(cluster)
		},
	}
)
