package controller

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
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
	}{
		{
			name:                     "exsisting namespace found",
			expectedName:             "my-namespace",
			targetLabelKey:           "label1",
			expectedLabelValue:       "value1",
			targetAnnotationsKey:     "remove-empty-ns-operator/will-removed",
			expectedAnnotationsValue: "True",
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
		}

		got, err := ctrl.GetNamepsaces()
		nsList := got.Items

		require.Equal(t, tc.expectedError, err, tc.name)

		require.Equal(t, nsList[0].ObjectMeta.Name, tc.expectedName)

		require.Equal(t, nsList[0].ObjectMeta.Labels[tc.targetLabelKey],
			tc.expectedLabelValue, tc.name,
		)

		require.Equal(t, tc.expectedAnnotationsValue,
			nsList[0].ObjectMeta.Annotations[tc.targetAnnotationsKey], tc.name,
		)
	}
}

func TestDeleteNamespace(t *testing.T) {
	testCases := []struct {
		name                string
		namespaces          []runtime.Object
		deletingNsName      string
		isNotFoundExpected  bool
		expectedAliveNsName string
	}{
		{
			name:                "deleting exist namespace",
			deletingNsName:      "delete-me",
			isNotFoundExpected:  true,
			expectedAliveNsName: "alive-ns",
			namespaces: []runtime.Object{
				&v1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "delete-me",
						Labels: map[string]string{
							"label1": "value1",
						},
					},
				},
				&v1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "alive-ns",
						Labels: map[string]string{
							"label2": "value2",
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

		ctrl.DeleteNamespace(tc.deletingNsName)

		_, err := ctrl.clientSet.CoreV1().Namespaces().Get(
			context.TODO(), tc.deletingNsName, metav1.GetOptions{},
		)
		// expect NotFoundError here
		require.Equal(t, tc.isNotFoundExpected, apierrs.IsNotFound(err), tc.name+": not found error")

		// expect that ns wasn't deleted
		existNs, err := ctrl.clientSet.CoreV1().Namespaces().Get(
			context.TODO(), tc.expectedAliveNsName, metav1.GetOptions{},
		)
		// ns must exist here
		require.Equal(t, false, apierrs.IsNotFound(err), tc.name+": unexpected not found error")
		// and gets without errors
		require.Equal(t, nil, err, tc.name+": unexpected error")

		require.Equal(t, tc.expectedAliveNsName, existNs.Name, tc.name)

	}
}
