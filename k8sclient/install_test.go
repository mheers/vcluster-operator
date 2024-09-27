package k8sclient

import (
	"testing"

	"github.com/mheers/vcluster-operator/config"
	"github.com/stretchr/testify/require"
)

func TestInstall(t *testing.T) {
	_, err := Init(&config.ServerConfig{
		K8sInCluster: false,
	})
	require.NoError(t, err)

	err = Install(&InstallOptions{})
	require.NoError(t, err)
}

func TestDumpYaml(t *testing.T) {
	_, err := Init(&config.ServerConfig{
		K8sInCluster: false,
	})
	require.NoError(t, err)

	err = Install(&InstallOptions{
		DumpYaml: true,
	})
	require.NoError(t, err)
}
