package internal

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

// getClientset Creates a new clientset for the kubernetes
func getClientset() (*kubernetes.Clientset, error) {
	var config *rest.Config

	// First try to use the in-cluster configuration
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fallback to kubeconfig
		var kubeconfig string
		if kc := os.Getenv("KUBECONFIG"); kc != "" {
			kubeconfig = kc
		} else {
			kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build config from kubeconfig path %s: %v", kubeconfig, err)
		}
	}

	// Create and store the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	return clientset, nil
}
