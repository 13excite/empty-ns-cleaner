package main

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("RUN")
	clientset, err := newClientSet(true)
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		ctx := context.Background()
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientset.CoreV1().Pods("test1").Get(ctx, "dnsutils", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod dnsutils not found in test1 namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found dnsutils pod in test1 namespace\n")
		}

		time.Sleep(10 * time.Second)
	}

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
