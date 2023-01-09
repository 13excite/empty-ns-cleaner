package controller

import (
	"context"
	"fmt"
	"log"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
	types "k8s.io/apimachinery/pkg/types"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	AddRemoveAnnotationValue = "True"
	DelRemoveAnnotationValue = "False"
)

type NSCleaner struct {
	kclient    *kubernetes.Clientset
	nsInformer cache.SharedIndexInformer // should i use cache for NS ?????

	dryRun bool
	ctx    context.Context
	stopCh <-chan struct{}
}

func NewNSCleaner(ctx context.Context, kclient *kubernetes.Clientset) *NSCleaner {
	return &NSCleaner{
		ctx:     ctx,
		kclient: kclient,
	}
}

func (c *NSCleaner) Run() {

}

func (c *NSCleaner) DeleteNamespace() {

}

func (c *NSCleaner) GetNamepsaces() (*v1.NamespaceList, error) {
	nsList, err := c.kclient.CoreV1().Namespaces().List(c.ctx, metav1.ListOptions{})
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
	_, err := c.kclient.CoreV1().Namespaces().Update(c.ctx, obj, metav1.UpdateOptions{})
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
// and adds remove-empty-ns-operator/will-removed=True
func (c *NSCleaner) patchWillRemovedAnnotations(name, annotationValue string) error {
	// default annotation value
	payload := fmt.Sprintf(
		`{"metadata": {"annotations": {"remove-empty-ns-operator/will-removed": "%s"}}}`,
		annotationValue,
	)
	// use MergePatchType here, because
	// the annotations field may not exist
	_, err := c.kclient.CoreV1().Namespaces().Patch(c.ctx, name, types.MergePatchType,
		[]byte(payload), metav1.PatchOptions{},
	)
	// notFoundErr is ok
	if ignoreNotFound(err) != nil {
		log.Print("ERROR IN PATCH")
		return err
	}
	return nil
}

// UpdateLablels uses official k8s way
// for updating labels
// todo state
func (c *NSCleaner) UpdateLablelsDISABLED() {
	// working with labels
	labels := map[string]string{
		"testlabel":  "value",
		"testlabel1": "value",
	}
	// shoud be pointer to obj
	//
	accessor, err := meta.Accessor(&struct{}{})
	if err != nil {
		log.Printf(err.Error())
	}

	objLabels := accessor.GetLabels()
	if objLabels == nil {
		objLabels = make(map[string]string)
	}

	for key, value := range labels {
		objLabels[key] = value
	}
	//fmt.Println(n)

	accessor.SetLabels(objLabels)
	// end labels blocks
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}
