package vclustermanagement

import (
	"testing"

	"github.com/mheers/vcluster-operator/config"
	"github.com/mheers/vcluster-operator/k8sclient"
	"github.com/stretchr/testify/require"
)

func TestKubeconfig(t *testing.T) {
	cfg := config.GetFakeServerConfig()
	_, err := k8sclient.Init(cfg)
	require.Nil(t, err)

	kubeconfig, err := Kubeconfig("test")
	require.Nil(t, err)
	require.NotEmpty(t, kubeconfig)
}
