package cmd

import (
	"fmt"

	"github.com/mheers/vcluster-operator/helpers"
	"github.com/spf13/cobra"
)

var (
	kubeconfigCmd = &cobra.Command{
		Use:   "kubeconfig [name]",
		Short: "gets the kubeconfig of a vcluster",
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
			kubeconfig, err := client.Kubeconfig(name)
			if err != nil {
				return err
			}

			fmt.Println(kubeconfig)
			return nil
		},
	}
)
