package controller

import (
	"context"
	"log"

	apierrs "k8s.io/apimachinery/pkg/api/errors"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}

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

// UpdateLabels updates given namespace
// Should use only for updating labels
// but also can be use for updating any fields
func (c *NSCleaner) UpdateLabels(obj *v1.Namespace) error {
	_, err := c.kclient.CoreV1().Namespaces().Update(c.ctx, obj, metav1.UpdateOptions{})
	if err != nil {
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
