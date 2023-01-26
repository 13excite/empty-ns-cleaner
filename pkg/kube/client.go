package kube

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// Clents contains 3 type of k8s clients
// DiscoveryClient works with server-supported API groups, versions and resources
// DynamicClient works with unstructured components
// ClientSet works with structed components
type Clients struct {
	ClientSet       *kubernetes.Clientset
	DiscoveryClient *discovery.DiscoveryClient
	DynamicClient   *dynamic.DynamicClient
}

func NewClients(outsideCluster bool) (*Clients, error) {
	config, err := newConfig(outsideCluster)
	if err != nil {
		return &Clients{}, err
	}

	dscvClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return &Clients{}, err
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return &Clients{}, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &Clients{}, err
	}

	return &Clients{
		DiscoveryClient: dscvClient,
		DynamicClient:   dynClient,
		ClientSet:       clientSet,
	}, nil
}
