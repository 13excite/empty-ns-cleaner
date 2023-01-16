package main

import (
	"context"
	"fmt"
	"github.com/13excite/empty-ns-cleaner/pkg/controller"
	//"github.com/13excite/empty-ns-cleaner/pkg/utils"
	//"k8s.io/apimachinery/pkg/api/errors"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

func main() {

	// don't mark this NS

	fmt.Println("RUN")
	clientset, err := newClientSet(true)
	if err != nil {
		log.Fatal(err.Error())
	}

	discClient, err := newDiscoveryClient(true)
	if err != nil {
		log.Fatal(err.Error())
	}

	dynClient, err := newDynamicClient(true)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()
	kContrl := controller.NewNSCleaner(ctx, clientset, discClient, dynClient)

	kContrl.Run(ctx)

}

func newClientSet(outsideCluster bool) (*kubernetes.Clientset, error) {
	kubeConfigPath := ""
	if outsideCluster {
		kubeConfigPath = filepath.Join(os.Getenv("KUBECONFIG"))
		fmt.Println(kubeConfigPath)
	}
	// if all args of BuildConfigFromFlags is ""
	// then rest.InClusterConfig() will be activeted
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	fmt.Println(config)

	return kubernetes.NewForConfig(config)
}

func newDynamicClient(outsideCluster bool) (*dynamic.DynamicClient, error) {
	kubeConfigPath := ""
	if outsideCluster {
		kubeConfigPath = filepath.Join(os.Getenv("KUBECONFIG"))
		fmt.Println(kubeConfigPath)
	}
	// if all args of BuildConfigFromFlags is ""
	// then rest.InClusterConfig() will be activeted
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return dynClient, nil
}

func newDiscoveryClient(outsideCluster bool) (*discovery.DiscoveryClient, error) {
	kubeConfigPath := ""
	if outsideCluster {
		kubeConfigPath = filepath.Join(os.Getenv("KUBECONFIG"))
		fmt.Println(kubeConfigPath)
	}
	// if all args of BuildConfigFromFlags is ""
	// then rest.InClusterConfig() will be activeted
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	dycoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	return dycoveryClient, nil
}
