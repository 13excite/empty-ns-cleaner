package controller

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// isEmpty checkes objects from the GV in a given namespaces
// and takes into account the ignored resources
func (c *NSCleaner) isEmpty(ns v1.Namespace,
	gvrList []schema.GroupVersionResource,
	workerLogValue string,
) bool {
GVR_LOOP:
	for _, gvr := range gvrList {
		objUnstruct, err := c.dynamicClient.Resource(gvr).Namespace(ns.Name).List(c.ctx, metav1.ListOptions{})
		if err != nil {
			if ignoreNotFound(err) != nil {
				c.logger.Errorw(`unexpected error of dynamic client.`, "GVR", gvr, "error", err.Error())
				continue GVR_LOOP
			}
			continue GVR_LOOP
		}
	OBJECT_LOOP:
		for _, obj := range objUnstruct.Items {
			if isIgnoredResouce(obj, gvr.Group, c.config.IgnoredResouces) {
				continue OBJECT_LOOP
			}
			c.logger.Debugw(
				"object exists in ns",
				"namespace", ns.Name,
				"kind", obj.Object["kind"],
				"name", obj.Object["metadata"].(map[string]interface{})["name"],
				"worker", workerLogValue,
			)
			return false
		}
	}
	return true
}

// DeleteNamespace deletes a namespace by name
func (c *NSCleaner) DeleteNamespace(name string) {
	propagation := metav1.DeletePropagationForeground
	if err := c.clientSet.CoreV1().Namespaces().Delete(c.ctx, name, metav1.DeleteOptions{
		PropagationPolicy: &propagation,
	}); ignoreNotFound(err) != nil {
		c.logger.Errorw("failed to delete ns '%s': %v", name, err)
		return
	}
	// TODO: add metrics
}

// GetNamepsaces returns a list of namespaces and an error
func (c *NSCleaner) GetNamepsaces() (*v1.NamespaceList, error) {
	nsList, err := c.clientSet.CoreV1().Namespaces().List(c.ctx, metav1.ListOptions{})
	if err != nil {
		c.logger.Fatalw(err.Error())
		return nil, err
	}
	return nsList, nil
}
