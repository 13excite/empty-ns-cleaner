package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/13excite/empty-ns-cleaner/pkg/config"
	"github.com/13excite/empty-ns-cleaner/pkg/utils"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
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
func NewNSCleaner(ctx context.Context, conf *config.Config,
	clientSet *kubernetes.Clientset,
	discoveryClient *discovery.DiscoveryClient,
	dynamicClient *dynamic.DynamicClient,
) *NSCleaner {
	return &NSCleaner{
		clientSet:       clientSet,
		discoveryClient: discoveryClient,
		dynamicClient:   dynamicClient,
		ctx:             ctx,
		config:          conf,
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

func (c *NSCleaner) Run(ctx context.Context) {
	for {

		namespaces, err := c.GetNamepsaces()
		gvRecouceList := c.GetApiRecources()

		if err != nil {
			panic(err.Error())
		}

		for _, n := range namespaces.Items {
			d := fmt.Sprintf("Found NS. Name: %s. Created: %v", n.Name, n.CreationTimestamp)

			if utils.IsContains(c.config.ProtectedNS, n.Name) {
				if c.config.DebugMode {
					log.Printf("NS %s is prodtected. Skiping....\n", n.Name)
				}
				continue
			}

			// DEBUG PRINT
			if c.isEmpty(n, gvRecouceList) {
				log.Printf("NS IS EMPTY: %s", n.Name)
			} else {
				log.Printf("NS IS NOT EMPTY: %s", n.Name)
			}
			// TODO: mark only empty namespaces
			// TODO: unmark annotations
			// working with labels
			// update labels
			if n.ObjectMeta.Annotations["remove-empty-ns-operator/will-removed"] != "True" {
				err := c.AddRemoveAnnotation(n.Name)
				if err != nil {
					log.Print(err)
				}
			} else {
				fmt.Printf("NS %s already marked as deleted\n", n.Name)
			}

			log.Print(d)
		}
		time.Sleep(time.Duration(c.config.RunEveeryMins) * time.Minute)
	}
}

// set as Public for testing
func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}
