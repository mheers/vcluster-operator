package vclustermanagement

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	clusters, err := List()
	require.Nil(t, err)
	require.NotEmpty(t, clusters)
}
