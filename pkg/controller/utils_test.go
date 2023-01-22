package controller

import (
	"testing"

	"github.com/13excite/empty-ns-cleaner/pkg/config"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestIsIgnoredResouce(t *testing.T) {
	testData := []struct {
		name             string
		want             bool
		object           unstructured.Unstructured
		ignoredResources []config.IgnoredResources
		apiGroup         string
	}{
		{
			name: "Full mask name matches",
			want: true,
			object: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "kube-root-ca.crt",
					},
					"kind": "ConfigMap",
				},
			},
			ignoredResources: []config.IgnoredResources{
				{
					NameMask: "kube-root-ca.crt",
					Kind:     "ConfigMap",
					APIGroup: "",
				},
			},
		},
		{
			name: "Regexp mask name matches",
			want: true,
			object: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "default-blabla",
					},
					"kind": "ServiceAccount",
				},
			},
			ignoredResources: []config.IgnoredResources{
				{
					NameMask: "^default-.*",
					Kind:     "ServiceAccount",
					APIGroup: "",
				},
			},
			apiGroup: "",
		},
		{
			name: "Regexp mask name doesn't matches",
			want: false,
			object: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "my-bad-pod",
					},
					"kind": "Pod",
				},
			},
			ignoredResources: []config.IgnoredResources{
				{
					NameMask: "^important-pod",
					Kind:     "Pod",
					APIGroup: "",
				},
			},
			apiGroup: "",
		},
		{
			name: "APIGroup doesn't matches",
			want: false,
			object: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "my-deployment",
					},
					"kind": "Deployment",
				},
			},
			ignoredResources: []config.IgnoredResources{
				{
					NameMask: "my-deployment",
					Kind:     "Deployment",
					APIGroup: "my.app.example.com/v1",
				},
			},
			apiGroup: "apps/v1",
		},
	}

	for _, tc := range testData {
		got := isIgnoredResouce(tc.object, tc.apiGroup, tc.ignoredResources, false)
		require.Equal(t, tc.want, got, tc.name)
	}
}
