package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/13excite/empty-ns-cleaner/pkg/config"
	"github.com/13excite/empty-ns-cleaner/pkg/kube"
	"github.com/13excite/empty-ns-cleaner/pkg/utils"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type NSCleaner struct {
	config *config.Config

	clientSet       kubernetes.Interface
	discoveryClient *discovery.DiscoveryClient
	dynamicClient   *dynamic.DynamicClient
	nsInformer      cache.SharedIndexInformer // should i use cache for NS ?????

	dryRun bool
	ctx    context.Context
	stopCh <-chan struct{}
}

// TODO: pass args via config struct
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
	}
}

// TODO: move logic to separate func
func (c *NSCleaner) Run() error {
	ticker := time.NewTicker(time.Duration(c.config.RunEveeryMins) * time.Minute)

	log.Printf("ns cleaner is starting....\n")

	for {
		select {
		case <-ticker.C:
			c.cleansingRunner()

		case <-c.ctx.Done():
			return nil
		}
	}
}

func (c *NSCleaner) cleansingRunner() {
	namespaces, err := c.GetNamepsaces()
	gvRecouceList := c.GetApiRecources()
	if err != nil {
		panic(err.Error())
	}

	for _, n := range namespaces.Items {
		if c.config.DebugMode {
			d := fmt.Sprintf("Found NS. Name: %s. Created: %v", n.Name, n.CreationTimestamp)
			log.Printf(d)
		}

		if utils.IsContains(c.config.ProtectedNS, n.Name) {
			if c.config.DebugMode {
				log.Printf("NS %s is prodtected. Skiping....\n", n.Name)
			}
			continue
		}

		shouldRemove := n.ObjectMeta.Annotations[CustomAnnotationName] == "True"

		if c.isEmpty(n, gvRecouceList) {
			// if ns is empty and has a deletion mark
			if shouldRemove {
				log.Printf("NS IS EMPTY AND HAS DELETION MARK: %s", n.Name)
				log.Printf("DELETING!!!!\n")
				// TODO: add a deletion method
				// if ns is empty and doesn't have a deletion mark
			} else {
				log.Printf("NS IS EMPTY AND DOESNT HAVE DELETION MARK: %s", n.Name)
				log.Printf("ADDDING DELETION MARK\n")
				c.AddWillRemoveAnnotation(n.Name)
			}
		} else {
			log.Printf("NS IS NOT EMPTY: %s", n.Name)
			// if ns isn't empty and has a deletion mark
			if shouldRemove {
				log.Printf("NS IS EMPTY AND HAS DELETION MARK: %s", n.Name)
				log.Printf("DELETING DELETION MARK")
				c.DeleteWillRemoveAnnotation(n.Name)
			}
		}
	}
}

func (c *NSCleaner) GetApiRecources() []schema.GroupVersionResource {
	// get resources list
	lists, err := c.discoveryClient.ServerPreferredResources()
	if err != nil {
		// TODO: log or return
		log.Printf(err.Error())
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
