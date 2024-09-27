package cmd

import (
	"fmt"

	"github.com/mheers/vcluster-operator/helpers"
	"github.com/spf13/cobra"
)

var (
	tokenCmd = &cobra.Command{
		Use:   "token [name]",
		Short: "token gets a jwt token to manage a vcluster",
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
			token, err := client.ClusterToken(name)
			if err != nil {
				return err
			}

			fmt.Println(token)

			return nil
		},
	}
)
