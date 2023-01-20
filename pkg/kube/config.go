package kube

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// newConfig returns a new k8s rest config
func newConfig(outsideCluster bool) (*rest.Config, error) {
	kubeConfigPath := ""
	if outsideCluster {
		kubeConfigPath = filepath.Join(os.Getenv("KUBECONFIG"))
	}
	// if all args of BuildConfigFromFlags is ""
	// then rest.InClusterConfig() will be activeted
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}
