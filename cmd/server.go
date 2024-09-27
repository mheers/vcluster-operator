package cmd

import (
	"github.com/mheers/vcluster-operator/config"
	"github.com/mheers/vcluster-operator/helpers"
	"github.com/mheers/vcluster-operator/k8sclient"
	"github.com/mheers/vcluster-operator/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "starts the server",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := config.GetConfig()

			_, err := k8sclient.Init(config)
			if err != nil {
				return err
			}

			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			// Create the server
			logrus.Info("Creating and starting the server")
			app := server.NewApplicaton(config)
			return app.Run()
		},
	}
)
