package k8sclient

import (
	"context"
	"flag"
	"path/filepath"
	"sync"

	"github.com/mheers/vcluster-operator/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	K8sClient *kubernetes.Clientset
	Ctx       = context.Background()
	once      sync.Once
)

// Init initializes a message queue client
func Init(appConfig *config.ServerConfig) (*kubernetes.Clientset, error) {
	var err error
	once.Do(func() {

		clusterConfig, err := getClusterConfig(appConfig)
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		clientset, err := kubernetes.NewForConfig(clusterConfig)
		if err != nil {
			panic(err.Error())
		}
		K8sClient = clientset
	})
	return K8sClient, err
}

func getClusterConfig(appConfig *config.ServerConfig) (*rest.Config, error) {
	if appConfig.K8sInCluster {
		return getInClusterConfig()
	}
	return getOutOfClusterConfig()
}

func getInClusterConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return config, err
}

func getOutOfClusterConfig() (*rest.Config, error) {
	var kubeconfig string
	var defaultPath string
	kcFlag := flag.Lookup("kubeconfig")
	if home := homedir.HomeDir(); home != "" {
		defaultPath = filepath.Join(home, ".kube", "config")
		if kcFlag == nil {
			flag.StringVar(&kubeconfig, "kubeconfig", defaultPath, "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = kcFlag.Value.String()
		}
	} else {
		defaultPath = ""
		if kcFlag == nil {
			flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
		} else {
			kubeconfig = kcFlag.Value.String()
		}
	}
	flag.Parse()

	if kubeconfig == "" {
		kubeconfig = defaultPath
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}
