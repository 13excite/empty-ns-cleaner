package controller

import (
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *NSCleaner) isEmpty(ns v1.Namespace, gvrList []schema.GroupVersionResource) bool {
GVR_LOOP:
	for _, gvr := range gvrList {
		objUnstruct, err := c.dynamicClient.Resource(gvr).Namespace(ns.Name).List(c.ctx, metav1.ListOptions{})
		if err != nil {
			if ignoreNotFound(err) != nil {
				log.Printf("unexpected error of dynamic client. GVP: %v \n", gvr)
				continue GVR_LOOP
			}
			continue GVR_LOOP
		}
	OBJECT_LOOP:
		for _, obj := range objUnstruct.Items {
			if isIgnoredResouce(obj, gvr.Group, c.config.IgnoredResouces, c.config.DebugMode) {
				continue OBJECT_LOOP
			}
			if c.config.DebugMode {
				log.Printf(
					"Name: %s KIND: %s is exist \n",
					obj.Object["metadata"].(map[string]interface{})["name"], obj.Object["kind"],
				)
			}
			return false
		}
	}
	return true
}

func (c *NSCleaner) DeleteNamespace(name string) {
	propagation := metav1.DeletePropagationForeground
	if err := c.clientSet.CoreV1().Namespaces().Delete(c.ctx, name, metav1.DeleteOptions{
		PropagationPolicy: &propagation,
	}); ignoreNotFound(err) != nil {
		log.Printf("failed to delete ns '%s': %v", name, err)
		return
	}
	// TODO: add metrics
}

func (c *NSCleaner) GetNamepsaces() (*v1.NamespaceList, error) {
	nsList, err := c.clientSet.CoreV1().Namespaces().List(c.ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return nsList, nil
}
