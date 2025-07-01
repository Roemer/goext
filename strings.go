package goext

import "strings"

// Checks if the string contains at least one of the substrings.
func StringContainsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
