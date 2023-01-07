package controller

import (
	"context"
	"log"

	apierrs "k8s.io/apimachinery/pkg/api/errors"

	v1 "k8s.io/api/core/v1"
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
