package controller

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

func TestAddWillRemoveAnnotation(t *testing.T) {
	testCases := []struct {
		name                     string
		namespaces               []runtime.Object
		nsName                   string
		targetAnnotationsKey     string
		expectedAnnotationsValue string
		expectedError            error
	}{
		{
			name:                     "adding will-removed annotations for existing ns",
			nsName:                   "my-namespace",
			targetAnnotationsKey:     CustomAnnotationName,
			expectedAnnotationsValue: "True",
			expectedError:            nil,
			namespaces: []runtime.Object{
				&v1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "my-namespace",
						Labels: map[string]string{
							"label1": "value1",
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
			logger:    zap.S().With("package", "testing"),
		}

		err := ctrl.AddWillRemoveAnnotation(tc.nsName)
		require.Equal(t, tc.expectedError, err, tc.name)

		got, err := ctrl.GetNamepsaces()

		require.Equal(t, tc.expectedError, err, tc.name+". Getting namespaces")

		require.Equal(t, tc.expectedAnnotationsValue,
			got.Items[0].ObjectMeta.Annotations[tc.targetAnnotationsKey], tc.name,
		)

	}
}

func TestDeleteRemoveAnnotation(t *testing.T) {
	testCases := []struct {
		name                     string
		namespaces               []runtime.Object
		nsName                   string
		targetAnnotationsKey     string
		expectedAnnotationsValue string
		expectedError            error
	}{
		{
			name:                     "deleting will-removed annotations for existing ns",
			nsName:                   "my-namespace",
			targetAnnotationsKey:     CustomAnnotationName,
			expectedAnnotationsValue: "False",
			expectedError:            nil,
			namespaces: []runtime.Object{
				&v1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "my-namespace",
						Labels: map[string]string{
							"label1": "value1",
						},
						Annotations: map[string]string{
							CustomAnnotationName: "True",
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
			logger:    zap.S().With("package", "ns-cleaner"),
		}

		err := ctrl.DeleteWillRemoveAnnotation(tc.nsName)
		require.Equal(t, tc.expectedError, err, tc.name)

		got, err := ctrl.GetNamepsaces()

		require.Equal(t, tc.expectedError, err, tc.name+". Getting namespaces")

		require.Equal(t, tc.expectedAnnotationsValue,
			got.Items[0].ObjectMeta.Annotations[tc.targetAnnotationsKey], tc.name,
		)
	}
}
