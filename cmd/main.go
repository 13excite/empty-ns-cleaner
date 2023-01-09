package main

import (
	"context"
	"fmt"
	"github.com/13excite/empty-ns-cleaner/pkg/controller"
	"github.com/13excite/empty-ns-cleaner/pkg/utils"
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

	// don't mark this NS
	protectedNS := []string{
		"default",
		"kube-public",
		"kube-system",
		"local-path-storage",
		"kube-node-lease",
	}

	fmt.Println("RUN")
	clientset, err := newClientSet(true)
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		ctx := context.Background()

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

		kContrl := controller.NewNSCleaner(ctx, clientset)

		namespaces, err := kContrl.GetNamepsaces()
		if err != nil {
			panic(err.Error())
		}

		for _, n := range namespaces.Items {
			d := fmt.Sprintf("Found NS. Name: %s. Created: %v", n.Name, n.CreationTimestamp)

			if utils.IsProtectedNs(protectedNS, n.Name) {
				fmt.Printf("NS %s is prodtected. Skiping....\n", n.Name)
				continue
			}

			// working with labels
			// update labels
			if n.ObjectMeta.Annotations["remove-empty-ns-operator/will-removed"] != "True" {
				err := kContrl.AddRemoveAnnotation(n.Name)
				if err != nil {
					log.Print(err)
				}
			} else {
				fmt.Printf("NS %s already marked as deleted\n", n.Name)
			}

			log.Print(d)
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
