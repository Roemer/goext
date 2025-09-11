package goext

import "strings"

// Checks if the string contains at least one of the substrings.
func StringContainsAny(s string, substrings ...string) bool {
	for _, sub := range substrings {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// Trims all occurrences of the given prefix from the start of the string.
func StringTrimAllPrefix(s, prefix string) string {
	for strings.HasPrefix(s, prefix) {
		s = s[len(prefix):]
	}
	return s
}

// Trims all occurrences of the given suffix from the end of the string.
func StringTrimAllSuffix(s, suffix string) string {
	for strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

// Trims all occurrences of newline characters from the end of the string.
func StringTrimNewlineSuffix(v string) string {
	return StringTrimAllSuffix(StringTrimAllSuffix(v, "\r\n"), "\n")
}

// Splits the string by new line characters, supporting both "\n" and "\r\n".
func StringSplitByNewLine(value string) []string {
	return strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), "\n")
}
