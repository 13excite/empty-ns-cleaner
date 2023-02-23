package utils

import (
	"sort"
)

// IsContains returns true is string exists
// in an ordered slice
func IsContains(s []string, e string) bool {
	i := sort.Search(len(s), func(i int) bool { return s[i] >= e })
	if i < len(s) && s[i] == e {
		return true
	} else {
		return false
	}
}
