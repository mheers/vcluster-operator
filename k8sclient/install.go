package k8sclient

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mheers/vcluster-operator/helpers"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

var deployDir = "./deploy"

type InstallOptions struct {
	NamespaceName          string
	DeploymentName         string
	ServiceAccountName     string
	SecretName             string
	ContainerName          string
	AppLabel               string
	ClusterRoleBindingName string
	ClusterRoleName        string
	Image                  string
	ImagePullPolicy        string
	AdminUser              string
	AdminPassword          string
	DumpYaml               bool
}

func (opts *InstallOptions) mergeWithDefaults() {
	var (
		namespaceName          = "vcluster-operator"
		deploymentName         = "vcluster-operator"
		serviceAccountName     = "vcluster-operator"
		secretName             = "vcluster-operator"
		containerName          = "vcluster-operator"
		appLabel               = "vcluster-operator"
		clusterRoleBindingName = "vcluster-operator-binding"
		clusterRoleName        = "vcluster-operator"
		image                  = "mheers/vcluster-operator:latest"
		imagePullPolicy        = "IfNotPresent"
		adminUser              = "admin"
		adminPassword          = ""
	)

	if opts.NamespaceName == "" {
		opts.NamespaceName = namespaceName
	}
	if opts.DeploymentName == "" {
		opts.DeploymentName = deploymentName
	}
	if opts.ServiceAccountName == "" {
		opts.ServiceAccountName = serviceAccountName
	}
	if opts.SecretName == "" {
		opts.SecretName = secretName
	}
	if opts.ContainerName == "" {
		opts.ContainerName = containerName
	}
	if opts.AppLabel == "" {
		opts.AppLabel = appLabel
	}
	if opts.ClusterRoleBindingName == "" {
		opts.ClusterRoleBindingName = clusterRoleBindingName
	}
	if opts.ClusterRoleName == "" {
		opts.ClusterRoleName = clusterRoleName
	}
	if opts.Image == "" {
		opts.Image = image
	}
	if opts.ImagePullPolicy == "" {
		opts.ImagePullPolicy = imagePullPolicy
	}
	if opts.AdminUser == "" {
		opts.AdminUser = adminUser
	}
	if opts.AdminPassword == "" {
		opts.AdminPassword = adminPassword
	}
}

func Install(opts *InstallOptions) error {
	opts.mergeWithDefaults()

	ctx := context.Background()

	if opts.DumpYaml {
		// mkdir -p ./deploy
		err := os.MkdirAll(deployDir, 0755)
		if err != nil {
			return err
		}
	}

	// Create namespace
	if opts.DumpYaml {
		fmt.Printf("Create namespace using kubectl create namespace %s\n", opts.NamespaceName)
	} else {
		found, err := NamespaceExists(ctx, opts.NamespaceName)
		if err != nil {
			return err
		}

		if !found {
			_, err = K8sClient.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: opts.NamespaceName,
				},
			}, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		}
	}

	// Create cluster role
	clusterRole := GetClusterRole(opts)
	if opts.DumpYaml {
		err := dumpYaml(clusterRole, "clusterrole")
		if err != nil {
			return err
		}
	} else {
		_, err := K8sClient.RbacV1().ClusterRoles().Create(ctx, clusterRole, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	// Create cluster role binding
	clusterRoleBinding := GetClusterRoleBinding(opts)
	if opts.DumpYaml {
		err := dumpYaml(clusterRoleBinding, "clusterrolebinding")
		if err != nil {
			return err
		}
	} else {
		_, err := K8sClient.RbacV1().ClusterRoleBindings().Create(ctx, clusterRoleBinding, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	// Create service account
	serviceAccount := GetServiceAccount(opts)
	if opts.DumpYaml {
		err := dumpYaml(serviceAccount, "serviceaccount")
		if err != nil {
			return err
		}
	} else {
		_, err := K8sClient.CoreV1().ServiceAccounts(opts.NamespaceName).Create(ctx, serviceAccount, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	// Create secret
	secret := GetSecret(opts)
	if opts.DumpYaml {
		err := dumpYaml(secret, "secret")
		if err != nil {
			return err
		}
	} else {
		_, err := K8sClient.CoreV1().Secrets(opts.NamespaceName).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	// Create deployment
	deployment := GetDeployment(opts)
	if opts.DumpYaml {
		err := dumpYaml(deployment, "deployment")
		if err != nil {
			return err
		}
	} else {
		_, err := K8sClient.AppsV1().Deployments(opts.NamespaceName).Create(ctx, deployment, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	// Create service
	service := GetService(opts)
	if opts.DumpYaml {
		err := dumpYaml(service, "service")
		if err != nil {
			return err
		}
	} else {
		_, err := K8sClient.CoreV1().Services(opts.NamespaceName).Create(ctx, service, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	// Create ingress
	ingress := GetIngress(opts)
	if opts.DumpYaml {
		err := dumpYaml(ingress, "ingress")
		if err != nil {
			return err
		}
	} else {
		_, err := K8sClient.NetworkingV1().Ingresses(opts.NamespaceName).Create(ctx, ingress, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	usedAdminUser := secret.Data["admin-user"]
	usedAdminPW := secret.Data["admin-password"]
	ingressUrl := ingress.Spec.Rules[0].Host

	if !opts.DumpYaml {
		fmt.Println("VCluster Operator installed successfully.")
	}
	fmt.Printf("\nLog in with user %s and password %s at %s:\n\n", usedAdminUser, usedAdminPW, ingressUrl)
	fmt.Printf("$ export VCLUSTER_OPERATOR_URL=http://%s:8080\n", ingressUrl) // TODO use https and the correct port
	fmt.Printf("$ vcluster-operator login --username %s --password %s\n", usedAdminUser, usedAdminPW)
	return nil
}

func GetDeployment(opts *InstallOptions) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.DeploymentName,
			Namespace: opts.NamespaceName,
			Labels: map[string]string{
				"app": opts.AppLabel,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": opts.AppLabel,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": opts.AppLabel,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: opts.ServiceAccountName,
					Containers: []corev1.Container{
						{
							Name:            opts.ContainerName,
							Image:           opts.Image,
							ImagePullPolicy: corev1.PullPolicy(opts.ImagePullPolicy),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
									Name:          "http",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "VCLUSTER_OPERATOR_K8S_INCLUSTER",
									Value: "true",
								},
								{
									Name:  "VCLUSTER_OPERATOR_PORT",
									Value: "8080",
								},
								{
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: opts.SecretName,
											},
											Key: "admin-user",
										},
									},
									Name: "VCLUSTER_OPERATOR_ADMIN_USER",
								},
								{
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: opts.SecretName,
											},
											Key: "admin-password",
										},
									},
									Name: "VCLUSTER_OPERATOR_ADMIN_PASSWORD",
								},
								{
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: opts.SecretName,
											},
											Key: "secret-key",
										},
									},
									Name: "VCLUSTER_OPERATOR_SECRET_KEY",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "kubeconfig",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "kubeconfig",
								},
							},
						},
					},
				},
			},
		},
	}

	return deployment
}

func GetService(opts *InstallOptions) *corev1.Service {
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.DeploymentName,
			Namespace: opts.NamespaceName,
			Labels: map[string]string{
				"app": opts.AppLabel,
			},
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": opts.AppLabel,
			},
			Ports: []corev1.ServicePort{
				{
					Port: 8080,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 8080,
					},
				},
			},
		},
	}

	return service
}

// GetServiceAccount
func GetServiceAccount(opts *InstallOptions) *corev1.ServiceAccount {
	serviceAccount := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.ServiceAccountName,
			Namespace: opts.NamespaceName,
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
		},
	}

	return serviceAccount
}

// GetClusterRoleBinding
func GetClusterRoleBinding(opts *InstallOptions) *rbacv1.ClusterRoleBinding {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.ClusterRoleBindingName,
			Namespace: opts.NamespaceName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      opts.ServiceAccountName,
				Namespace: opts.NamespaceName,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole",
			Name: opts.ClusterRoleName,
		},
	}

	return clusterRoleBinding
}

// GetClusterRole
func GetClusterRole(opts *InstallOptions) *rbacv1.ClusterRole {
	clusterRole := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.ClusterRoleName,
			Namespace: opts.NamespaceName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
		},
	}

	return clusterRole
}

// Ingress
func GetIngress(opts *InstallOptions) *networkingv1.Ingress {
	ingress := &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.DeploymentName,
			Namespace: opts.NamespaceName,
			Labels: map[string]string{
				"app": opts.AppLabel,
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "vcluster.mheers.dev",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: "/",
									PathType: func() *networkingv1.PathType {
										pt := networkingv1.PathTypePrefix
										return &pt
									}(),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: opts.DeploymentName,
											Port: networkingv1.ServiceBackendPort{
												Number: 8080,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return ingress
}

// GetSecret
func GetSecret(opts *InstallOptions) *corev1.Secret {
	adminUser := opts.AdminUser
	adminPassword := opts.AdminPassword

	if adminUser == "" {
		adminUser = "admin"
	}
	if adminPassword == "" {
		adminPassword = GetRandomSecretKey()
	}

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.SecretName,
			Namespace: opts.NamespaceName,
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
		},
		Data: map[string][]byte{
			"secret-key":     []byte(GetRandomSecretKey()),
			"admin-user":     []byte(adminUser),
			"admin-password": []byte(adminPassword),
		},
	}
	return secret
}

func dumpYaml(obj interface{}, name string) error {
	y, err := helpers.MarshalViaJSONToYAML(obj)
	if err != nil {
		return err
	}
	dst := filepath.Join(deployDir, name+".yaml")
	err = os.WriteFile(dst, y, 0644)
	if err != nil {
		return err
	}
	return nil
}
