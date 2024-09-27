package vclustermanagement

import (
	"github.com/gin-gonic/gin"
	"github.com/loft-sh/vcluster/cmd/vclusterctl/cmd"
	"github.com/loft-sh/vcluster/cmd/vclusterctl/log"
)

// Create creates a new vcluster.  It returns the name and the kubeconfig of the new vcluster.
func Create(name string) (*cmd.VCluster, error) {
	root, err := cmd.BuildRoot(log.GetInstance())
	if err != nil {
		return nil, err
	}

	// see https://github.com/loft-sh/vcluster/blob/main/cmd/vclusterctl/cmd/create.go#L93
	// "kube-config-context-name", "", "If set, will override the context name of the generated virtual cluster kube config with this name"
	// "chart-version", upgrade.GetVersion(), "The virtual cluster chart version to use (e.g. v0.9.1)"
	// "chart-name", "vcluster", "The virtual cluster chart name to use"
	// "chart-repo", LoftChartRepo, "The virtual cluster chart repo to use"
	// "local-chart-dir", "", "The virtual cluster local chart dir to use"
	// "k3s-image", "", "DEPRECATED: use --extra-values instead"
	// "distro", "k3s", fmt.Sprintf("Kubernetes distro to use for the virtual cluster. Allowed distros: %s", strings.Join(AllowedDistros, ", "))
	// "release-values", "", "DEPRECATED: use --extra-values instead"
	// "kubernetes-version", "", "The kubernetes version to use (e.g. v1.20). Patch versions are not supported"
	// "extra-values", "f", []string{}, "Path where to load extra helm values from"
	// "create-namespace", true, "If true the namespace will be created if it does not exist"
	// "disable-ingress-sync", false, "If true the virtual cluster will not sync any ingresses"
	// "update-current", true, "If true updates the current kube config"
	// "create-cluster-role", false, "DEPRECATED: cluster role is now automatically created if it is required by one of the resource syncers that are enabled by the .sync.RESOURCE.enabled=true helm value, which is set in a file that is passed via --extra-values argument."
	// "expose", false, "If true will create a load balancer service to expose the vcluster endpoint"
	// "expose-local", true, "If true and a local Kubernetes distro is detected, will deploy vcluster with a NodePort service"

	// "connect", true, "If true will run vcluster connect directly after the vcluster was created"
	// "upgrade", false, "If true will try to upgrade the vcluster instead of failing if it already exists"
	// "isolate", false, "If true vcluster and its workloads will run in an isolated environment"

	root.SetArgs([]string{
		"create",
		name,
		"--connect=false",
		"--namespace",
		getNamespaceName(name),
	})
	err = root.Execute()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func CreateHandler(c *gin.Context) {
	// claims := jwt.ExtractClaims(c)
	// user, _ := c.Get(config.IdentityKey)
	cluster, err := Create(c.Param("name"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cluster)
}
