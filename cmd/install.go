package cmd

import (
	"github.com/mheers/vcluster-operator/config"
	"github.com/mheers/vcluster-operator/helpers"
	"github.com/mheers/vcluster-operator/k8sclient"
	"github.com/spf13/cobra"
)

var (
	image           string
	imagePullPolicy string
	dumpYaml        bool
	installCmd      = &cobra.Command{
		Use:   "install [name]",
		Short: "installs the vcluster-operator",
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

			opts := &k8sclient.InstallOptions{
				Image:           image,
				ImagePullPolicy: imagePullPolicy,
				AdminUser:       username,
				AdminPassword:   password,
				DumpYaml:        dumpYaml,
			}
			err = k8sclient.Install(opts)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	installCmd.Flags().StringVarP(&username, "admin-username", "U", "", "username")
	installCmd.Flags().StringVarP(&password, "admin-password", "P", "", "password")
	installCmd.Flags().StringVarP(&image, "image", "i", "mheers/vcluster-operator:latest", "image")
	installCmd.Flags().StringVarP(&imagePullPolicy, "image-pull-policy", "p", "Always", "image pull policy")
	installCmd.Flags().BoolVarP(&dumpYaml, "dump-yaml", "d", false, "dump yaml to stdout instead of applying it")
}
