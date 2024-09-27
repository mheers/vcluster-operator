package vclustermanagement

import (
	"fmt"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mheers/vcluster-operator/k8sclient"
)

// Kubeconfig lists all available vclusters.  It returns a list of vcluster names.
func Kubeconfig(name string) ([]byte, error) {
	nameI := fmt.Sprintf("vc-%s", name)
	namespaceI := fmt.Sprintf("vcluster-%s", name)

	c := k8sclient.K8sClient
	ctx := k8sclient.Ctx
	secret, err := c.CoreV1().Secrets(namespaceI).Get(ctx, nameI, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	config := secret.Data["config"]
	return config, nil
}

func KubeconfigHandler(c *gin.Context) {
	name := c.Params.ByName("name")
	kubeconfig, err := Kubeconfig(name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.String(200, string(kubeconfig))
}
