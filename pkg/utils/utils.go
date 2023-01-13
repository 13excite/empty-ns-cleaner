package utils

import (
	"sort"
)

// IsProtectedNs returns true is namespace exists in
// protected namespaces list. BinarySearch uses here
// because a cluster can containe a couple hundred namespaces
func IsProtectedNs(protectedNS []string, namespaceName string) bool {
	i := sort.Search(len(protectedNS), func(i int) bool { return protectedNS[i] >= namespaceName })

	if i < len(protectedNS) && protectedNS[i] == namespaceName {
		return true
	} else {
		return false
	}
}

// ContainsVerb returns true if string exists in slice
// Simple iterations uses here, because because
// verbSlice contains as usual less then 10 items
func ContainsVerb(s []string, v string) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
}
