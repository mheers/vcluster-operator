package k8sclient

import (
	"testing"

	"github.com/mheers/vcluster-operator/config"
	"github.com/stretchr/testify/require"
)

func TestUninstall(t *testing.T) {
	_, err := Init(&config.ServerConfig{
		K8sInCluster: false,
	})
	require.NoError(t, err)

	err = Uninstall(&InstallOptions{})
	require.NoError(t, err)
}
