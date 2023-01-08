package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsProtectedNs(t *testing.T) {
	testData := []struct {
		name       string
		potectedNS []string
		inputNS    string
		want       bool
	}{
		{
			name: "Protected namespace",
			potectedNS: []string{
				"my-important-ns1",
				"my-important-ns2",
				"my-important-ns3",
				"my-important-ns4",
			},
			inputNS: "my-important-ns3",
			want:    true,
		},
		{
			name: "Non protected namespace",
			potectedNS: []string{
				"my-important-ns1",
				"my-important-ns2",
				"my-important-ns3",
				"my-important-ns4",
			},
			inputNS: "my-unimportant-ns",
			want:    false,
		},
	}

	for _, tc := range testData {
		got := IsProtectedNs(tc.potectedNS, tc.inputNS)
		require.Equal(t, tc.want, got, tc.name)
	}

}
