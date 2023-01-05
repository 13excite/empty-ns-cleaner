package controller


import (
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"

)

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}

type NSCleaner struct {
	DryRun bool
	kclient     *kubernetes.Clientset
}