package controller

import (
	"log"
	"regexp"

	"github.com/13excite/empty-ns-cleaner/pkg/config"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// isIgnoredResouce returns true if object is
// exist in config's IgnoredResources slice
func isIgnoredResouce(obj unstructured.Unstructured, APIGroup string,
	ignoredResources []config.IgnoredResources, debugMode bool) bool {
	for _, ignoreResource := range ignoredResources {
		// type casting a resource's name to string
		objName := obj.Object["metadata"].(map[string]interface{})["name"].(string)
		matchIsIgnored, err := regexp.MatchString(ignoreResource.NameMask, objName)
		// if regexp match was failed,
		// write log and return false
		if err != nil {
			log.Printf("couldn't match string with error: %v", err)
			return false
		}
		// if full match between a config and resource,
		// then returns true
		if matchIsIgnored &&
			(ignoreResource.Kind == obj.Object["kind"]) &&
			(ignoreResource.APIGroup == APIGroup) {
			if debugMode {
				log.Printf("IGNORRED RESOURCES FROM CONFIG. KIND:%s NAME: %s \n",
					ignoreResource.Kind, objName)
			}
			return true
		}
	}
	return false
}
