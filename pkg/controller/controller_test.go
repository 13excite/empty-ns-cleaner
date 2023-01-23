package controller

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/require"
)

func TestGetNamepsaces(t *testing.T) {
	testCases := []struct {
		name                     string
		namespaces               []runtime.Object
		expectedName             string
		targetLabelKey           string
		expectedLabelValue       string
		targetAnnotationsKey     string
		expectedAnnotationsValue string
		expectedError            error
		expectSuccess            bool
	}{
		{
			name:                     "exsisting namespace found",
			expectedName:             "my-namespace",
			targetLabelKey:           "label1",
			expectedLabelValue:       "value1",
			targetAnnotationsKey:     "remove-empty-ns-operator/will-removed",
			expectedAnnotationsValue: "True",
			expectedError:            nil,
			expectSuccess:            true,
			namespaces: []runtime.Object{
				&v1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "my-namespace",
						Labels: map[string]string{
							"label1": "value1",
						},
						Annotations: map[string]string{
							"remove-empty-ns-operator/will-removed": "True",
						},
					},
				},
			},
		},
		{
			name:                     "exsisting namespace found",
			expectedName:             "my-namespace",
			targetLabelKey:           "label1",
			expectedLabelValue:       "value1",
			targetAnnotationsKey:     "remove-empty-ns-operator/will-removed",
			expectedAnnotationsValue: "True",
			expectedError:            nil,
			expectSuccess:            true,
			namespaces: []runtime.Object{
				&v1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "my-namespace",
						Labels: map[string]string{
							"label1": "value1",
						},
						Annotations: map[string]string{
							"remove-empty-ns-operator/will-removed": "True",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		fakeClientset := fake.NewSimpleClientset(tc.namespaces...)
		ctrl := &NSCleaner{
			clientSet: fakeClientset,
		}

		got, err := ctrl.GetNamepsaces()
		nsList := got.Items

		require.Equal(t, tc.expectedError, err, tc.name)

		require.Equal(t, nsList[0].ObjectMeta.Name, tc.expectedName)

		require.Equal(t, nsList[0].ObjectMeta.Labels[tc.targetLabelKey],
			tc.expectedLabelValue, tc.name,
		)

		require.Equal(t, nsList[0].ObjectMeta.Annotations[tc.targetAnnotationsKey],
			tc.expectedAnnotationsValue, tc.name,
		)
	}
}
