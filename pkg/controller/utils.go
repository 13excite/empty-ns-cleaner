package controller

import (
	"regexp"

	"github.com/13excite/empty-ns-cleaner/pkg/config"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
)

// isIgnoredResouce returns true if object is
// exist in config's IgnoredResources slice
func isIgnoredResouce(obj unstructured.Unstructured,
	APIGroup string,
	ignoredResources []config.IgnoredResources,
) bool {
	for _, ignoreResource := range ignoredResources {
		// type casting a resource's name to string
		objName := obj.Object["metadata"].(map[string]interface{})["name"].(string)
		matchIsIgnored, err := regexp.MatchString(ignoreResource.NameMask, objName)
		// if regexp match was failed,
		// write log and return false
		if err != nil {
			zap.S().Errorw("couldn't match string", "error", err)
			return false
		}
		// if full match between a config and resource,
		// then returns true
		if matchIsIgnored &&
			(ignoreResource.Kind == obj.Object["kind"]) &&
			(ignoreResource.APIGroup == APIGroup) {
			zap.S().Debugw("ignorred resource from config", "kind", ignoreResource.Kind, "name", objName)
			return true
		}
	}
	return false
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}
