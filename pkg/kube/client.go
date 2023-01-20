package kube

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// NewDiscoveryClient returns a client for
// discover server-supported API groups, versions and resources
func NewDiscoveryClient(outsideCluster bool) (*discovery.DiscoveryClient, error) {
	config, err := newConfig(outsideCluster)
	if err != nil {
		return nil, err
	}
	return discovery.NewDiscoveryClientForConfig(config)
}

// NewDynamicClient returns a new dynamic client
func NewDynamicClient(outsideCluster bool) (*dynamic.DynamicClient, error) {
	config, err := newConfig(outsideCluster)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

// NewClientSet returns a Clientset which contains the clients for groups
func NewClientSet(outsideCluster bool) (*kubernetes.Clientset, error) {
	config, err := newConfig(outsideCluster)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
