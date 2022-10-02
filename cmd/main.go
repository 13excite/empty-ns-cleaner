package main

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("RUN")
	clientset, err := newClientSet()
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx := context.Background()
}

func newClientSet() (*kubernetes.Clientset, error) {

	kubeConfigPath := filepath.Join(os.Getenv("KUBECONFIG"))
	fmt.Println(kubeConfigPath)

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	fmt.Println(config)

	return kubernetes.NewForConfig(config)
}
