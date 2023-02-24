package main

import (
	"context"
	"flag"
	"log"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/13excite/empty-ns-cleaner/pkg/config"
	"github.com/13excite/empty-ns-cleaner/pkg/controller"
	"github.com/13excite/empty-ns-cleaner/pkg/kube"
	"github.com/13excite/empty-ns-cleaner/pkg/utils"
)

func main() {

	isOutsideCluster := flag.Bool("outside", true, "is service running outside of k8s")
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	config.C.Defaults()
	if *configPath != "" {
		if err := config.C.ReadConfig(*configPath); err != nil {
			log.Printf("could not read config file with error: %v\n", err)
		}
	}
	config.InitLogger(&config.C)
	logger := zap.S().With("package", "cmd")

	kubeClients, err := kube.NewClients(*isOutsideCluster)
	if err != nil {
		logger.Fatalw(err.Error())
	}

	// TODO: TRY TO SPLIT IT
	// SHOULD TO CREATE SEPARATE STRUCT SERVICE(or similar name) ????
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context, cancel context.CancelFunc) {
		defer cancel()
		utils.WaitForShutdown(ctx)
	}(ctx, cancel)

	group, ctx := errgroup.WithContext(ctx)
	kContrl := controller.NewNSCleaner(ctx, &config.C, kubeClients)

	group.Go(func() error {
		return kContrl.Run()
	})

	err = group.Wait()
	if err != nil {
		logger.Panicf("service runs error")
	}
}
