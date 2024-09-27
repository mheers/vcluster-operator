package vclustermanagement

import (
	"github.com/gin-gonic/gin"
	"github.com/loft-sh/vcluster/cmd/vclusterctl/cmd/find"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Get lists all available vclusters.  It returns a list of vcluster names.
func Get(name string) (*find.VCluster, error) {
	context := ""
	namespace := metav1.NamespaceAll
	vClusters, err := find.GetVCluster(context, name, namespace)
	if err != nil {
		return nil, err
	}

	return vClusters, nil
}

func GetHandler(c *gin.Context) {
	// claims := jwt.ExtractClaims(c)
	// user, _ := c.Get(config.IdentityKey)
	cluster, err := Get(c.Param("name"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cluster)
}
