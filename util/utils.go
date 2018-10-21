package util

import "strings"

// HumanizeBool outputs a human friendly string representation of a bool
func HumanizeBool(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

// StringIsEmpty checks to see if a string is empty
func StringIsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
