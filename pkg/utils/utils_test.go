package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsContains(t *testing.T) {
	testData := []struct {
		name  string
		input []string
		exist string
		want  bool
	}{
		{
			name: "Protected namespace",
			input: []string{
				"my-important-ns1",
				"my-important-ns2",
				"my-important-ns3",
				"my-important-ns4",
			},
			exist: "my-important-ns3",
			want:  true,
		},
		{
			name: "Non protected namespace",
			input: []string{
				"my-important-ns1",
				"my-important-ns2",
				"my-important-ns3",
				"my-important-ns4",
			},
			exist: "my-unimportant-ns",
			want:  false,
		},
		{
			name:  "String exists",
			input: []string{"create", "delete", "deletecollection", "get", "list"},
			exist: "get",
			want:  true,
		},
		{
			name:  "String doesn't exist",
			input: []string{"create", "delete", "deletecollection", "list", "watch"},
			exist: "get",
			want:  false,
		},
	}

	for _, tc := range testData {
		got := IsContains(tc.input, tc.exist)
		require.Equal(t, tc.want, got, tc.name)
	}
}
