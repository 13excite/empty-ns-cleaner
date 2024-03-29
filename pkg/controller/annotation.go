package controller

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

const (
	// CustomAnnotationName is a default annotations name for this service
	CustomAnnotationName     = "remove-empty-ns-operator/will-removed"
	AddRemoveAnnotationValue = "True"
	DelRemoveAnnotationValue = "False"
)

// AddRemoveAnnotation removes deletion mark
// and add remove-empty-ns-operator/will-removed=False
func (c *NSCleaner) DeleteWillRemoveAnnotation(name string) error {
	return c.patchWillRemovedAnnotations(name, DelRemoveAnnotationValue)
}

// AddRemoveAnnotation adds deletion
// mark remove-empty-ns-operator/will-removed=True
func (c *NSCleaner) AddWillRemoveAnnotation(name string) error {
	return c.patchWillRemovedAnnotations(name, AddRemoveAnnotationValue)
}

// PatchWillRemovedAnnotations patches annotations of namespace
// and adds remove-empty-ns-operator/will-removed=${annotationValue}
func (c *NSCleaner) patchWillRemovedAnnotations(name, annotationValue string) error {
	// default annotation value
	payload := fmt.Sprintf(
		`{"metadata": {"annotations": {"%s": "%s"}}}`,
		CustomAnnotationName,
		annotationValue,
	)
	// use MergePatchType here, because
	// the annotations field may not exist
	_, err := c.clientSet.CoreV1().Namespaces().Patch(c.ctx, name, types.MergePatchType,
		[]byte(payload), metav1.PatchOptions{},
	)
	// notFoundErr is ok
	if ignoreNotFound(err) != nil {
		c.logger.Errorw("unexpected error in running patch annotations", "error", err)
		return err
	}
	return nil
}
