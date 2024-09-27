package k8sclient

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Uninstall(opts *InstallOptions) error {
	opts.mergeWithDefaults()

	ctx := context.Background()

	// we delete the namespace
	namespaceExists, err := NamespaceExists(ctx, opts.NamespaceName)
	if err != nil {
		return err
	}
	if namespaceExists {
		fmt.Println("Deleting namespace", opts.NamespaceName)
		err := K8sClient.CoreV1().Namespaces().Delete(ctx, opts.NamespaceName, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	// we delete the clusterrolebinding
	clusterRoleBindingExists, err := ClusterRoleBindingExists(ctx, opts.ClusterRoleBindingName)
	if err != nil {
		return err
	}
	if clusterRoleBindingExists {
		fmt.Println("Deleting clusterrolebinding", opts.ClusterRoleBindingName)
		err := K8sClient.RbacV1().ClusterRoleBindings().Delete(ctx, opts.ClusterRoleBindingName, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	// we delete the clusterrole
	clusterRoleExists, err := ClusterRoleExists(ctx, opts.ClusterRoleName)
	if err != nil {
		return err
	}
	if clusterRoleExists {
		fmt.Println("Deleting clusterrole", opts.ClusterRoleName)
		err := K8sClient.RbacV1().ClusterRoles().Delete(ctx, opts.ClusterRoleName, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
