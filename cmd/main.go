package main

import (
	"context"
	"flag"
	"log"

	"github.com/13excite/empty-ns-cleaner/pkg/config"
	"github.com/13excite/empty-ns-cleaner/pkg/controller"
	"github.com/13excite/empty-ns-cleaner/pkg/kube"
	//"github.com/13excite/empty-ns-cleaner/pkg/utils"
	//"k8s.io/apimachinery/pkg/api/errors"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {

	isOutsideCluster := flag.Bool("outside", true, "is service running outside of k8s")
	flag.Parse()

	config.C.Defaults()

	log.Println("RUN")
	clientset, err := kube.NewClientSet(*isOutsideCluster)
	if err != nil {
		log.Fatal(err.Error())
	}

	discClient, err := kube.NewDiscoveryClient(*isOutsideCluster)
	if err != nil {
		log.Fatal(err.Error())
	}

	dynClient, err := kube.NewDynamicClient(*isOutsideCluster)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()
	kContrl := controller.NewNSCleaner(ctx, &config.C,
		clientset, discClient, dynClient,
	)

	kContrl.Run(ctx)

}
