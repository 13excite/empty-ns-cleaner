package controller

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/13excite/empty-ns-cleaner/pkg/config"
	"github.com/13excite/empty-ns-cleaner/pkg/kube"
	"github.com/13excite/empty-ns-cleaner/pkg/utils"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type NSCleaner struct {
	config *config.Config
	// TODO: update logic for adding worker number to the logs
	logger *zap.SugaredLogger

	clientSet       kubernetes.Interface
	discoveryClient *discovery.DiscoveryClient
	dynamicClient   *dynamic.DynamicClient
	nsInformer      cache.SharedIndexInformer // should i use cache for NS ?????

	dryRun bool
	ctx    context.Context
	stopCh <-chan struct{}
}

func NewNSCleaner(
	ctx context.Context,
	conf *config.Config,
	kubeCleints *kube.Clients,
) *NSCleaner {
	return &NSCleaner{
		ctx:             ctx,
		clientSet:       kubeCleints.ClientSet,
		discoveryClient: kubeCleints.DiscoveryClient,
		dynamicClient:   kubeCleints.DynamicClient,
		config:          conf,
		// should to switch from sugar to a base logger?
		logger: zap.S().With("package", "ns-cleaner"),
	}
}

// TODO: move logic to separate func
func (c *NSCleaner) Run() error {
	ticker := time.NewTicker(time.Duration(c.config.RunEveeryMins) * time.Minute)
	c.logger.Infow("ns cleaner is starting....")

	for {
		select {
		case <-ticker.C:
			c.cleaningRunner()
		case <-c.ctx.Done():
			ticker.Stop()
			return nil
		}
	}
}

func (c *NSCleaner) cleaningRunner() {
	namespaces, err := c.GetNamepsaces()
	gvRecouceList := c.getApiRecources()
	if err != nil {
		panic(err.Error())
	}
	// all available CPUs
	runtime.GOMAXPROCS(0)

	wg := &sync.WaitGroup{}
	workerInput := make(chan *v1.Namespace, c.config.NumWorkers)
	defer close(workerInput)

	for i := 0; i < c.config.NumWorkers; i++ {
		go c.cleaningWorker(workerInput, wg, gvRecouceList, i)
	}

	for _, n := range namespaces.Items {
		wg.Add(1)
		workerInput <- &n
	}
	wg.Wait()
	c.logger.Infow("all workers finished")
}

func (c *NSCleaner) getApiRecources() []schema.GroupVersionResource {
	// get resources list
	lists, err := c.discoveryClient.ServerPreferredResources()
	if err != nil {
		// TODO: log or return
		c.logger.Errorw("error of getting api recources", "error", err)
	}
	// result recources
	resources := []schema.GroupVersionResource{}
	for _, list := range lists {
		if len(list.APIResources) == 0 {
			continue
		}
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			continue
		}
	LOOP_API_RESOURCES:
		for _, resource := range list.APIResources {
			if len(resource.Verbs) == 0 {
				continue LOOP_API_RESOURCES
			}
			// skip recources without "get" method
			if !utils.IsContains(resource.Verbs, "get") {
				continue LOOP_API_RESOURCES
			}
			// skip Events
			if resource.Name == "events" {
				continue LOOP_API_RESOURCES
			}
			// skip cluster-wide recources, like
			// clusterRoles and etc
			if !resource.Namespaced {
				continue LOOP_API_RESOURCES
			}

			resources = append(resources, schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.String(),
				Resource: resource.Name,
			})
		}
	}
	return resources
}
