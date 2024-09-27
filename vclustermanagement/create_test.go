package vclustermanagement

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	clusters, err := Create("demo")
	require.Nil(t, err)
	require.NotEmpty(t, clusters)
}
