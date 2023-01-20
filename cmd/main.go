package main

import (
	"context"
	"fmt"
	"github.com/13excite/empty-ns-cleaner/pkg/config"
	"github.com/13excite/empty-ns-cleaner/pkg/controller"
	"github.com/13excite/empty-ns-cleaner/pkg/kube"

	//"github.com/13excite/empty-ns-cleaner/pkg/utils"
	//"k8s.io/apimachinery/pkg/api/errors"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func main() {

	config.C.Defaults()

	fmt.Println("RUN")
	clientset, err := kube.NewClientSet(true)
	if err != nil {
		log.Fatal(err.Error())
	}

	discClient, err := kube.NewDiscoveryClient(true)
	if err != nil {
		log.Fatal(err.Error())
	}

	dynClient, err := kube.NewDynamicClient(true)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()
	kContrl := controller.NewNSCleaner(ctx, &config.C,
		clientset, discClient, dynClient,
	)

	kContrl.Run(ctx)

}
