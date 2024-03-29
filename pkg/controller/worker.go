package controller

import (
	"runtime"
	"strconv"
	"sync"

	"github.com/13excite/empty-ns-cleaner/pkg/utils"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *NSCleaner) cleaningWorker(
	inNamespace <-chan *v1.Namespace,
	wg *sync.WaitGroup,
	gvRecouceList []schema.GroupVersionResource,
	workerNum int,
) {
	workerLogValue := "worker-" + strconv.Itoa(workerNum)
	for n := range inNamespace {
		// default values for logger
		fields := []interface{}{
			"namespace", n.Name, "worker", workerLogValue,
		}
		c.logger.Debugw("found NS", fields...)

		if utils.IsContains(c.config.ProtectedNS, n.Name) {
			c.logger.Debugw("protected ns was skipped ", fields...)
			// also paste done wg here
			wg.Done()
			continue
		}

		shouldRemove := n.ObjectMeta.Annotations[CustomAnnotationName] == "True"

		if c.isEmpty(*n, gvRecouceList, workerLogValue) {
			// if ns is empty and has a deletion mark
			if shouldRemove {
				c.logger.Debugw("NS is empty and has deletion mark", fields...)
				c.logger.Infow("deleting ns", fields...)
				c.DeleteNamespace(n.Name)
				// if ns is empty and doesn't have a deletion mark
			} else {
				c.logger.Debugw("NS is empty and doesn't have deletion mark", fields...)
				c.logger.Infow("adding deletion mark", fields...)
				err := c.AddWillRemoveAnnotation(n.Name)
				if err != nil {
					errMsgFields := append(fields, "error", err.Error())
					c.logger.Errorw("annotation error", errMsgFields...)
				}
			}
		} else {
			c.logger.Debugw("NS is not empty", fields...)
			// if ns isn't empty and has a deletion mark
			if shouldRemove {
				c.logger.Debugw("NS is NOT empty and has deletion mark", fields...)
				c.logger.Infow("deleting deletion mark", fields...)
				err := c.DeleteWillRemoveAnnotation(n.Name)
				if err != nil {
					errMsgFields := append(fields, "error", err.Error())
					c.logger.Errorw("deletion error", errMsgFields...)
				}
			}
		}
		wg.Done()
		runtime.Gosched()
	}
}
