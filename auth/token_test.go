package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	SecretKey := "secret"
	AdminUser := "admin"
	AdminPassword := "admin"
	authMiddleware, err := GetAuthMiddleware(SecretKey, AdminUser, AdminPassword)
	require.NoError(t, err)
	require.NotNil(t, authMiddleware)

	token, err := authMiddleware.Token("test")
	require.NoError(t, err)
	require.NotEmpty(t, token)
}
