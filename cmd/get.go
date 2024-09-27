package cmd

import (
	"github.com/loft-sh/vcluster/cmd/vclusterctl/cmd/find"
	"github.com/mheers/vcluster-operator/helpers"
	"github.com/spf13/cobra"
)

var (
	getCmd = &cobra.Command{
		Use:   "get [name]",
		Short: "gets a vcluster",
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
			cluster, err := client.Get(name)
			if err != nil {
				return err
			}

			return renderCluster(cluster)
		},
	}
)

func renderCluster(cluster *find.VCluster) error {
	if OutputFormatFlag == "table" {
		renderListTable([]find.VCluster{*cluster})
	}
	if OutputFormatFlag == "json" {
		err := helpers.PrintJSON(cluster)
		if err != nil {
			return err
		}
	}
	if OutputFormatFlag == "yaml" {
		err := helpers.PrintYAML(cluster)
		if err != nil {
			return err
		}
	}
	if OutputFormatFlag == "csv" {
		err := helpers.PrintCSV(cluster)
		if err != nil {
			return err
		}
	}
	return nil
}
