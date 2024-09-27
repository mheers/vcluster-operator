package k8sclient

import (
	"context"
	"testing"

	"github.com/mheers/vcluster-operator/config"
	"github.com/stretchr/testify/require"
)

func TestNamespaceExists(t *testing.T) {
	_, err := Init(&config.ServerConfig{
		K8sInCluster: false,
	})
	require.NoError(t, err)

	exists, err := NamespaceExists(context.Background(), "default")
	require.NoError(t, err)
	require.True(t, exists)

	notExists, err := NamespaceExists(context.Background(), "slejwpoytuirtj")
	require.NoError(t, err)
	require.False(t, notExists)
}

func TestClusterRoleBindingExists(t *testing.T) {
	_, err := Init(&config.ServerConfig{
		K8sInCluster: false,
	})
	require.NoError(t, err)

	exists, err := ClusterRoleBindingExists(context.Background(), "cluster-admin")
	require.NoError(t, err)
	require.True(t, exists)

	notExists, err := ClusterRoleBindingExists(context.Background(), "slejwpoytuirtj")
	require.NoError(t, err)
	require.False(t, notExists)
}

func TestClusterRoleExists(t *testing.T) {
	_, err := Init(&config.ServerConfig{
		K8sInCluster: false,
	})
	require.NoError(t, err)

	exists, err := ClusterRoleExists(context.Background(), "admin")
	require.NoError(t, err)
	require.True(t, exists)

	notExists, err := ClusterRoleExists(context.Background(), "slejwpoytuirtj")
	require.NoError(t, err)
	require.False(t, notExists)
}

func TestGetRandomSecretKey(t *testing.T) {
	key := GetRandomSecretKey()
	require.NotEmpty(t, key)
}
