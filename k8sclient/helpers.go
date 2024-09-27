package k8sclient

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NamespaceExists(ctx context.Context, namespaceName string) (bool, error) {
	ns, err := K8sClient.CoreV1().Namespaces().Get(ctx, namespaceName, metav1.GetOptions{})
	notFound := ns.ObjectMeta.Name != namespaceName

	// check if err is notfound
	status, ok := err.(*k8sErrors.StatusError)
	if ok || errors.As(err, &status) {
		if status.ErrStatus.Reason == metav1.StatusReasonNotFound {
			notFound = true
		}
	}
	if !notFound && err != nil {
		return false, err
	}
	return !notFound, nil
}

func ClusterRoleBindingExists(ctx context.Context, name string) (bool, error) {
	crb, err := K8sClient.RbacV1().ClusterRoleBindings().Get(ctx, name, metav1.GetOptions{})
	notFound := crb.ObjectMeta.Name != name
	// check if err is notfound
	status, ok := err.(*k8sErrors.StatusError)
	if ok || errors.As(err, &status) {
		if status.ErrStatus.Reason == metav1.StatusReasonNotFound {
			notFound = true
		}
	}
	if !notFound && err != nil {
		return false, err
	}
	return !notFound, nil
}

func ClusterRoleExists(ctx context.Context, name string) (bool, error) {
	cr, err := K8sClient.RbacV1().ClusterRoles().Get(ctx, name, metav1.GetOptions{})
	notFound := cr.ObjectMeta.Name != name
	// check if err is notfound
	status, ok := err.(*k8sErrors.StatusError)
	if ok || errors.As(err, &status) {
		if status.ErrStatus.Reason == metav1.StatusReasonNotFound {
			notFound = true
		}
	}
	if !notFound && err != nil {
		return false, err
	}
	return !notFound, nil
}

// RandomSecretKey
func GetRandomSecretKey() string {
	secretKey := make([]byte, 32)
	_, err := rand.Read(secretKey)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(secretKey)
}
