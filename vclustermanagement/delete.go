package vclustermanagement

import (
	"github.com/gin-gonic/gin"
	"github.com/loft-sh/vcluster/cmd/vclusterctl/cmd"
	"github.com/loft-sh/vcluster/cmd/vclusterctl/flags"
)

// Delete creates a new vcluster.  It returns the name and the kubeconfig of the new vcluster.
func Delete(name string) (*cmd.VCluster, error) {
	globalFlags := &flags.GlobalFlags{}
	vClusterCmd := cmd.NewDeleteCmd(globalFlags)

	// see https://github.com/loft-sh/vcluster/blob/main/cmd/vclusterctl/cmd/delete.go#L66
	// "keep-pvc", false, "If enabled, vcluster will not delete the persistent volume claim of the vcluster"
	// "delete-namespace", false, "If enabled, vcluster will delete the namespace of the vcluster"
	// "auto-delete-namespace", true, "If enabled, vcluster will delete the namespace of the vcluster if it was created by vclusterctl"
	vClusterCmd.SetArgs([]string{
		name,
		"--namespace",
		getNamespaceName(name),
	})
	vClusterCmd.Execute()
	return nil, nil
}

func DeleteHandler(c *gin.Context) {
	// claims := jwt.ExtractClaims(c)
	// user, _ := c.Get(config.IdentityKey)
	cluster, err := Delete(c.Param("name"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cluster)
}
