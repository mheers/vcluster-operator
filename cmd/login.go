package cmd

import (
	"errors"
	"os"

	"github.com/mheers/vcluster-operator/client"
	"github.com/mheers/vcluster-operator/helpers"
	"github.com/spf13/cobra"
)

var (
	url      string
	username string
	password string
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "login into the server",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			client, err := initClient(false)
			if err != nil {
				return err
			}
			err = client.Login(username, password)
			if err != nil {
				return err
			}
			cmd.Println("Successfully logged in")
			return nil
		},
	}
)

func init() {
	// clientCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&url, "url", "u", "", "url of the server")
	loginCmd.Flags().StringVarP(&username, "username", "U", "", "username")
	loginCmd.Flags().StringVarP(&password, "password", "P", "", "password")
}

func initClient(loadConfig bool) (*client.Client, error) {
	url, err := getURL()
	if err != nil {
		return nil, err
	}
	client := client.NewClient(url)
	if loadConfig {
		err = client.LoadConfig()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func getURL() (string, error) {
	if url != "" {
		return url, nil
	}
	urlFromEnv := os.Getenv("VCLUSTER_OPERATOR_URL")
	if urlFromEnv != "" {
		return urlFromEnv, nil
	}
	return "", errors.New("no url provided")
}
