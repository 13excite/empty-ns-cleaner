package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/13excite/empty-ns-cleaner/pkg/utils"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	v1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

const (
	AddRemoveAnnotationValue = "True"
	DelRemoveAnnotationValue = "False"
)

type NSCleaner struct {
	kClient         *kubernetes.Clientset
	discoveryClient *discovery.DiscoveryClient
	dynamicClient   *dynamic.DynamicClient
	nsInformer      cache.SharedIndexInformer // should i use cache for NS ?????
	dryRun          bool
	ctx             context.Context
	stopCh          <-chan struct{}
}

// TODO: pass args via config struct
func NewNSCleaner(ctx context.Context, kclient *kubernetes.Clientset,
	discoveryClient *discovery.DiscoveryClient,
	dynamicClient *dynamic.DynamicClient,
) *NSCleaner {
	return &NSCleaner{
		ctx:             ctx,
		kClient:         kclient,
		discoveryClient: discoveryClient,
		dynamicClient:   dynamicClient,
	}
}

// type groupResource struct {
// 	APIGroup        string
// 	APIGroupVersion string
// 	APIResource     metav1.APIResource
// }

func (c *NSCleaner) GetApiRecources() []schema.GroupVersionResource {
	// get resources list
	lists, err := c.discoveryClient.ServerPreferredResources()
	if err != nil {
		// TODO: log or return
		log.Println(err)
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
	protectedNS := []string{
		"default",
		"kube-public",
		"kube-system",
		"local-path-storage",
		"kube-node-lease",
	}
	for {

		// Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		// _, err = clientset.CoreV1().Pods("test1").Get(ctx, "dnsutils", metav1.GetOptions{})
		// if errors.IsNotFound(err) {
		// 	fmt.Printf("Pod dnsutils not found in test1 namespace\n")
		// } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		// 	fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		// } else if err != nil {
		// 	panic(err.Error())
		// } else {
		// 	fmt.Printf("Found dnsutils pod in test1 namespace\n")
		// }

		namespaces, err := c.GetNamepsaces()
		gvRecouceList := c.GetApiRecources()

		if err != nil {
			panic(err.Error())
		}

		for _, n := range namespaces.Items {
			d := fmt.Sprintf("Found NS. Name: %s. Created: %v", n.Name, n.CreationTimestamp)

			if utils.IsContains(protectedNS, n.Name) {
				fmt.Printf("NS %s is prodtected. Skiping....\n", n.Name)
				continue
			}

		GVR_LOOP:
			for _, gvr := range gvRecouceList {
				objUnstruct, err := c.dynamicClient.Resource(gvr).Namespace(n.Name).List(ctx, metav1.ListOptions{})
				if err != nil {
					//log.Print("GVR: ", gvr)
					if ignoreNotFound(err) != nil {
						log.Print("ERRRORRRR!!!!!! ", gvr)
						continue GVR_LOOP
					}
					continue GVR_LOOP
				}
				for _, obj := range objUnstruct.Items {
					fmt.Println("NAMESPACES: ", n.Name, "GROUP-RESOURCE", gvr.Group, gvr.Resource)
					fmt.Printf(
						"Name: %s KIND: %s \n",
						obj.Object["metadata"].(map[string]interface{})["name"], obj.Object["kind"],
					)
				}

			}

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

		time.Sleep(10 * time.Second)
	}

}

func (c *NSCleaner) DeleteNamespace() {

}

func (c *NSCleaner) GetNamepsaces() (*v1.NamespaceList, error) {
	nsList, err := c.kClient.CoreV1().Namespaces().List(c.ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return nsList, nil
}

// Updates given namespace
// Should use only for updating labels
// but also can be use for updating any fields
func (c *NSCleaner) update(obj *v1.Namespace) error {
	_, err := c.kClient.CoreV1().Namespaces().Update(c.ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// AddRemoveAnnotation removes deletion mark
// and add remove-empty-ns-operator/will-removed=False
func (c *NSCleaner) DeleteRemoveAnnotation(name string) error {
	return c.patchWillRemovedAnnotations(name, DelRemoveAnnotationValue)
}

// AddRemoveAnnotation adds deletion
// mark remove-empty-ns-operator/will-removed=True
func (c *NSCleaner) AddRemoveAnnotation(name string) error {
	return c.patchWillRemovedAnnotations(name, AddRemoveAnnotationValue)
}

// PatchWillRemovedAnnotations patches annotations of namespace
// and adds remove-empty-ns-operator/will-removed=${annotationValue}
func (c *NSCleaner) patchWillRemovedAnnotations(name, annotationValue string) error {
	// default annotation value
	payload := fmt.Sprintf(
		`{"metadata": {"annotations": {"remove-empty-ns-operator/will-removed": "%s"}}}`,
		annotationValue,
	)
	// use MergePatchType here, because
	// the annotations field may not exist
	_, err := c.kClient.CoreV1().Namespaces().Patch(c.ctx, name, types.MergePatchType,
		[]byte(payload), metav1.PatchOptions{},
	)
	// notFoundErr is ok
	if ignoreNotFound(err) != nil {
		log.Print("ERROR IN PATCH")
		return err
	}
	return nil
}

func (c *NSCleaner) isEmpty() {
	//dynamic.NewForConfigAndClient()
}

// set as Public for testing
func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}
