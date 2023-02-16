package controller

import (
	"context"
	"time"

	"github.com/13excite/empty-ns-cleaner/pkg/config"
	"github.com/13excite/empty-ns-cleaner/pkg/kube"
	"github.com/13excite/empty-ns-cleaner/pkg/utils"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type NSCleaner struct {
	config *config.Config
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
		logger:          zap.S().With("package", "ns-cleaner"),
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

	for _, n := range namespaces.Items {
		c.logger.Debugw("found NS", "name", n.Name, "created", n.CreationTimestamp)

		if utils.IsContains(c.config.ProtectedNS, n.Name) {
			c.logger.Debugw("protected ns was skipped ", "name", n.Name)
			continue
		}

		shouldRemove := n.ObjectMeta.Annotations[CustomAnnotationName] == "True"

		if c.isEmpty(n, gvRecouceList) {
			// if ns is empty and has a deletion mark
			if shouldRemove {
				c.logger.Infow("NS is empty and has deletion mark", "name", n.Name)
				c.logger.Debugw("deleting ns", "name", n.Name)
				c.DeleteNamespace(n.Name)
				// if ns is empty and doesn't have a deletion mark
			} else {
				c.logger.Infow("NS is empty and doesn't have deletion mark", "name", n.Name)
				c.logger.Infow("adding deletion mark", "name", n.Name)
				c.AddWillRemoveAnnotation(n.Name)
			}
		} else {
			c.logger.Infow("NS is not empty", "name", n.Name)
			// if ns isn't empty and has a deletion mark
			if shouldRemove {
				c.logger.Infow("NS is NOT empty and has deletion mark", "name", n.Name)
				c.logger.Infow("deleting deletion mark", "name", n.Name)
				c.DeleteWillRemoveAnnotation(n.Name)
			}
		}
	}
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
