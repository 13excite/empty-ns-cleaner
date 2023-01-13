package utils

import (
	"sort"
)

func IsProtectedNs(protectedNS []string, namespaceName string) bool {
	i := sort.Search(len(protectedNS), func(i int) bool { return protectedNS[i] >= namespaceName })

	if i < len(protectedNS) && protectedNS[i] == namespaceName {
		return true
	} else {
		return false
	}
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
