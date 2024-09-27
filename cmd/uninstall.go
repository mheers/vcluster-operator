package cmd

import (
	"github.com/mheers/vcluster-operator/config"
	"github.com/mheers/vcluster-operator/helpers"
	"github.com/mheers/vcluster-operator/k8sclient"
	"github.com/spf13/cobra"
)

var (
	uninstallCmd = &cobra.Command{
		Use:   "uninstall [name]",
		Short: "uninstalls the vcluster-operator",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			config := &config.ServerConfig{
				K8sInCluster: false,
			}

			_, err := k8sclient.Init(config)
			if err != nil {
				return err
			}

			// TODO: be able to uninstall a specific vcluster-operator (pram namespace)
			opts := &k8sclient.InstallOptions{}
			err = k8sclient.Uninstall(opts)
			if err != nil {
				return err
			}

			return nil
		},
	}
)
