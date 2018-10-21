package util

import (
	"fmt"
	"strings"

	"github.com/ryanuber/columnize"
)

// Columnize prints data in neat data columns
func Columnize(data []string) {
	config := columnize.DefaultConfig()
	config.Glue = "    "

	output := columnize.Format(data, config)
	fmt.Println()
	fmt.Println(output)
}

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
