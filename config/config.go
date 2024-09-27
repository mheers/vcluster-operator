package config

import (
	"context"
	"os"

	"github.com/imdario/mergo"
	"github.com/mheers/vcluster-operator/helpers"
	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
)

var IdentityKey = "id"

// ServerConfig describes the config
type ServerConfig struct {
	K8sInCluster  bool   `env:"VCLUSTER_OPERATOR_K8S_INCLUSTER"`
	Port          int    `env:"VCLUSTER_OPERATOR_PORT,default=8080"`
	AdminUser     string `env:"VCLUSTER_OPERATOR_ADMIN_USER"`
	AdminPassword string `env:"VCLUSTER_OPERATOR_ADMIN_PASSWORD"`
	SecretKey     string `env:"VCLUSTER_OPERATOR_SECRET_KEY"`
}

// OverlayConfigWithEnv overlays the config with values from the env
func (cfg *ServerConfig) OverlayConfigWithEnv() error {
	ctx := context.Background()
	overlayCfg := &ServerConfig{}
	err := envconfig.Process(ctx, overlayCfg)
	if err != nil {
		return err
	}

	err = mergo.Merge(cfg, overlayCfg, mergo.WithOverride)
	if err != nil {
		return err
	}
	return nil
}

// GetFakeServerConfig creates a config for testing purposes only
func GetFakeServerConfig() *ServerConfig {
	cfg := &ServerConfig{
		K8sInCluster: false,
	}

	err := cfg.OverlayConfigWithEnv()
	if err != nil {
		return nil
	}

	logLevel := os.Getenv("LOGLEVEL")
	if logLevel != "" {
		helpers.SetLogLevel(logLevel)
	}

	return cfg
}

var configInstance *ServerConfig

func GetConfig() *ServerConfig {
	if configInstance == nil {
		configInstance = &ServerConfig{}
		err := configInstance.OverlayConfigWithEnv()
		if err != nil {
			logrus.Fatalf("config: %s", err)
		}
	}
	return configInstance
}
