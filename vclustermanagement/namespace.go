package vclustermanagement

import "fmt"

func getNamespaceName(clusterName string) string {
	namespacePrefix := "vcluster"
	return fmt.Sprintf("%s-%s", namespacePrefix, clusterName)
}
