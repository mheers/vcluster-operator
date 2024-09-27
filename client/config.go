package client

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
	"k8s.io/client-go/util/homedir"
)

var configFolder = ".vcluster-operator"

type ClientConfig struct {
	URL         string    `yaml:"url"`
	Token       string    `yaml:"token"`
	TokenExpire time.Time `yaml:"tokenExpire"`
}

func createConfigFolder() (string, error) {
	if home := homedir.HomeDir(); home != "" {
		cf := filepath.Join(home, configFolder)
		return cf, os.MkdirAll(cf, 0700)
	}
	return "", errors.New("could not find home directory to store the client config")
}

func createConfigFilePath() (string, error) {
	absConfigFolder, err := createConfigFolder()
	if err != nil {
		return "", err
	}
	cfp := filepath.Join(absConfigFolder, "config.yaml")
	// create the file if it does not exist
	_, err = os.Stat(cfp)
	if os.IsNotExist(err) {
		_, err := os.Create(cfp)
		if err != nil {
			return "", err
		}
	}
	return cfp, nil
}

func getConfigFilePath() (string, error) {
	return createConfigFilePath()
}

func saveAll(configs []*ClientConfig) error {
	cfp, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(configs)
	if err != nil {
		return err
	}
	return os.WriteFile(cfp, data, 0600)
}

func readAll() ([]*ClientConfig, error) {
	cfp, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(cfp)
	if err != nil {
		return nil, err
	}
	configs := []*ClientConfig{}
	err = yaml.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (c *ClientConfig) Load() error {
	configs, err := readAll()
	if err != nil {
		return err
	}
	for _, config := range configs {
		if config.URL == c.URL {
			c.Token = config.Token
			c.TokenExpire = config.TokenExpire
			return nil
		}
	}
	return errors.New("could not find config for url " + c.URL)
}

func (c *ClientConfig) Save() error {
	configs, err := readAll()
	if err != nil {
		return err
	}
	for i, config := range configs {
		if config.URL == c.URL {
			configs[i] = c
			return saveAll(configs)
		}
	}
	configs = append(configs, c)
	return saveAll(configs)
}

func (c *ClientConfig) Delete() error {
	configs, err := readAll()
	if err != nil {
		return err
	}
	for i, config := range configs {
		if config.URL == c.URL {
			configs = append(configs[:i], configs[i+1:]...)
			return saveAll(configs)
		}
	}
	return nil
}
