package auth

import (
	"github.com/gin-gonic/gin"
)

// Token creates and returns a jwt token to manage a vcluster
func (m *Middleware) Token(name string) (string, error) {

	token, _, err := m.TokenGenerator(map[string]string{})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *Middleware) TokenHandler(c *gin.Context) {
	// claims := jwt.ExtractClaims(c)
	// user, _ := c.Get(config.IdentityKey)

	cluster, err := m.Token(c.Param("name"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cluster)
}
