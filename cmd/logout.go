package cmd

import (
	"github.com/mheers/vcluster-operator/helpers"
	"github.com/spf13/cobra"
)

var (
	logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "logout off the server",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			client, err := initClient(false)
			if err != nil {
				return err
			}
			err = client.Logout()
			if err != nil {
				return err
			}
			cmd.Println("Successfully logged out")
			return nil
		},
	}
)
