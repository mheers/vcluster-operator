package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/loft-sh/vcluster/cmd/vclusterctl/cmd/find"
	"github.com/mheers/vcluster-operator/helpers"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "list all vclusters",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			client, err := initClient(true)
			if err != nil {
				return err
			}
			clusters, err := client.List()
			if err != nil {
				return err
			}

			return renderClusters(clusters)
		},
	}
)

func renderClusters(clusters []find.VCluster) error {
	if OutputFormatFlag == "table" {
		renderListTable(clusters)
	}
	if OutputFormatFlag == "json" {
		err := helpers.PrintJSON(clusters)
		if err != nil {
			return err
		}
	}
	if OutputFormatFlag == "yaml" {
		err := helpers.PrintYAML(clusters)
		if err != nil {
			return err
		}
	}
	if OutputFormatFlag == "csv" {
		err := helpers.PrintCSV(clusters)
		if err != nil {
			return err
		}
	}
	return nil
}

func renderListTable(clusters []find.VCluster) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Namespace", "Status", "Created"})
	for _, cluster := range clusters {
		t.AppendRow(table.Row{cluster.Name, cluster.Namespace, cluster.Status, cluster.Created.Time})
		t.AppendSeparator()
	}
	t.Render()
}
